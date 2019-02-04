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

build: create-bulid-dir
	CGO_ENABLED=1 go build -o build/${PLUGIN_NAME}-${GOOS}-${GOARCH}${GOARM}.so -buildmode=plugin ${PLUGIN_ENTRY}

check: check-go-mod

.PHONY: build