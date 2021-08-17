BINARY=argocd-env-plugin
VERSION=0.0.1
OS_ARCH=darwin_amd64

default: build

build:
	go build -o ${BINARY} .

install: build
