ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all:
	go build -o ../bin/libxcbindings-darwin-amd64.so -buildmode=c-shared main.go
	mv ../bin/libxcbindings-darwin-amd64.h ../bin/libxcbindings.h
	docker run -v $(ROOT_DIR):/source --rm golang:1.10 /source/build-golang-linux
	mv  *.so ../bin
	rm *.h
