.PHONY: deps

image-proxy: image_proxy.go deps
	GOPATH=`pwd`/vendor go build -o $@ $<

deps:
	GOPATH=`pwd`/vendor go get "github.com/codegangsta/negroni" "github.com/nfnt/resize" "github.com/golang/groupcache/singleflight"
