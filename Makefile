BUILDDIR=./build
GOTIFY_VERSION=master
PLUGIN_NAME=myplugin
PLUGIN_ENTRY=plugin.go

download-tools:
	GO111MODULE=off go get -u github.com/gotify/plugin-api/cmd/gomod-cap

create-bulid-dir:
	mkdir -p ${BUILDDIR} || true

check-go-mod: create-bulid-dir
	wget -LO ${BUILDDIR}/gotify-server.mod https://raw.githubusercontent.com/gotify/server/${GOTIFY_VERSION}/go.mod
	gomod-cap -from ${BUILDDIR}/gotify-server.mod -to go.mod -check=true
	rm ${BUILDDIR}/gotify-server.mod || true

build-linux-amd64:
	 docker run --rm -v "$$PWD/.:/proj" -w /proj gotify/build:v1-linux-amd64 go build -a -installsuffix cgo -ldflags "-w -s" -buildmode=plugin -o build/${PLUGIN_NAME}-linux-amd64.so /proj

build-linux-arm-7:
	 docker run --rm -v "$$PWD/.:/proj" -w /proj gotify/build:v1-linux-arm-7 go build -a -installsuffix cgo -ldflags "-w -s" -buildmode=plugin -o build/${PLUGIN_NAME}-linux-arm-7.so /proj

build-linux-arm64:
	 docker run --rm -v "$$PWD/.:/proj" -w /proj gotify/build:v1-linux-arm64 go build -a -installsuffix cgo -ldflags "-w -s" -buildmode=plugin -o build/${PLUGIN_NAME}-linux-arm64.so /proj

build: build-linux-arm-7 build-linux-amd64 build-linux-arm64

check: check-go-mod

.PHONY: build
