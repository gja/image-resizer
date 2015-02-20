.PHONY: deps

main: deps
	GOPATH=`pwd`/vendor go build main.go

deps:
	GOPATH=`pwd`/vendor go get "github.com/codegangsta/negroni" "github.com/nfnt/resize" "github.com/golang/groupcache/singleflight"
