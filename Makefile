.PHONY: deps

main: deps
	go build main.go

deps:
	go get "github.com/codegangsta/negroni" "github.com/nfnt/resize"
