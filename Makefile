.PHONY: deps

SOURCES=image_proxy.go

image-proxy: $(SOURCES) deps
	GOPATH=`pwd`/vendor go build -o $@ $<

linux/image-proxy: $(SOURCES) deps
	GOOS=linux GOARCH=amd64 GOPATH=`pwd`/vendor go build -o $@ $<

deps:
	GOPATH=`pwd`/vendor go get "github.com/codegangsta/negroni" "github.com/nfnt/resize" "github.com/golang/groupcache/singleflight"
