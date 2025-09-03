package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotify/plugin-api/v2"
	"github.com/gotify/plugin-api/v2/generated/protobuf"
	"github.com/gotify/plugin-api/v2/transport"
	myplugin "github.com/gotify/plugin-template"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	httpTimeout = 15 * time.Second
)

type PluginServer struct {
	shutdown     chan struct{}
	shutdownOnce *sync.Once
	grpcServer   *grpc.Server
	ginEngine    *gin.Engine
	http.Server
	protobuf.UnimplementedPluginServer
	protobuf.UnimplementedDisplayerServer
	protobuf.UnimplementedConfigurerServer
}

func NewPluginServer(cli *plugin.PluginCli) *PluginServer {
	rootCAs := x509.NewCertPool()

	certificateChain, err := cli.Kex(myplugin.PluginInfo.ModulePath, rootCAs)
	if err != nil {
		log.Fatal(err)
	}
	tlsConfig := &tls.Config{
		Certificates: certificateChain,
		RootCAs:      rootCAs,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    rootCAs,
	}

	server := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             15 * time.Second,
		PermitWithoutStream: true,
	}), grpc.Creds(credentials.NewTLS(tlsConfig)))

	ginEngine := gin.Default()

	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetHTTP2(true)

	self := &PluginServer{
		grpcServer:   server,
		ginEngine:    ginEngine,
		shutdown:     make(chan struct{}),
		shutdownOnce: &sync.Once{},
	}

	self.Server = http.Server{
		Handler:           self,
		TLSConfig:         tlsConfig,
		Protocols:         protocols,
		ReadTimeout:       httpTimeout,
		ReadHeaderTimeout: httpTimeout,
		WriteTimeout:      httpTimeout,
	}

	return self
}

func (s *PluginServer) GetPluginInfo(ctx context.Context, req *emptypb.Empty) (*protobuf.Info, error) {
	return &myplugin.PluginInfo, nil
}

func (s *PluginServer) GracefulShutdown(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	s.shutdownOnce.Do(func() {
		close(s.shutdown)
	})
	return &emptypb.Empty{}, nil
}

func (s *PluginServer) RunUserInstance(req *protobuf.UserInstanceRequest, stream protobuf.Plugin_RunUserInstanceServer) error {
	for i := range myplugin.Capabilities {
		if err := stream.Send(&protobuf.InstanceUpdate{
			Update: &protobuf.InstanceUpdate_Capable{
				Capable: myplugin.Capabilities[i],
			},
		}); err != nil {
			return err
		}
	}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := stream.Send(&protobuf.InstanceUpdate{
				Update: &protobuf.InstanceUpdate_Ping{
					Ping: new(emptypb.Empty),
				},
			}); err != nil {
				return err
			}
		case <-s.shutdown:
			return nil
		}
	}
}

func (h *PluginServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil {
		http.Error(w, "Must use TLS", http.StatusUpgradeRequired)
		return
	}

	pluginRpcHostName := transport.BuildPluginTLSName(transport.PurposePluginRPC, myplugin.PluginInfo.ModulePath)

	if r.TLS.ServerName == pluginRpcHostName {
		if r.ProtoMajor != 2 {
			http.Error(w, "Must use HTTP/2", http.StatusHTTPVersionNotSupported)
			return
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			http.Error(w, "Must use application/grpc content type", http.StatusUnsupportedMediaType)
			return
		}
		h.grpcServer.ServeHTTP(w, r)

		return
	}

	pluginWebhookHostName := transport.BuildPluginTLSName(transport.PurposePluginWebhook, myplugin.PluginInfo.ModulePath)
	if r.TLS.ServerName == pluginWebhookHostName {
		h.ginEngine.ServeHTTP(w, r)
		return
	}

	http.Error(w, "Virtual host not found", http.StatusNotFound)
}

func main() {
	cli, err := plugin.ParsePluginCli(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	if cli.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	server := NewPluginServer(cli)

	listener, address, err := transport.NewListener()
	if err != nil {
		log.Fatal(err)
	}
	if err := server.ServeTLS(listener, "", ""); err != nil {
		log.Fatal(err)
	}

	_ = address
}
