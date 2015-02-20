package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"net/http"
	"strings"
	"os"
	"io"
)

type CachingResponseWriter struct {
	cacheLength int
	responseWriter http.ResponseWriter
	status int
}

func (c CachingResponseWriter) Header() http.Header {
	return c.responseWriter.Header()
}

func (c CachingResponseWriter) Write(bytes []byte) (int, error) {
	return c.responseWriter.Write(bytes)
}

func (c CachingResponseWriter) WriteHeader(status int) {
	if(status == 200) {
		c.Header().Set("Cache-Control", fmt.Sprintf("no-transform,public,max-age=%d,s-maxage=%d", c.cacheLength, c.cacheLength))
	}
	c.responseWriter.WriteHeader(status)
}

func cacheFor(time int) negroni.Handler {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		next(CachingResponseWriter{cacheLength: time, status: 200, responseWriter: rw}, r)
	})
}

func getExtension(filename string) string {
	x := strings.Split(filename, ".")
	return x[len(x) - 1]
}

func loadObject(path string) {
	parts := strings.Split(path, "/")

	if (len(parts) < 4) {
		return
	}
	props, key, filename := parts[1], strings.Join(parts[2:len(parts) - 1], "/"), parts[len(parts) - 1]
	outputDirectory := strings.Join(parts[1:len(parts) - 1], "/")
	extension := getExtension(filename)

	resp, err := http.Get("http://d1pcxoetpnw26i.cloudfront.net/thequint/2015-02/115f1834-306f-48ef-9024-7145e51e2cbe/manjhi.jpg-large")
	defer resp.Body.Close()
	if (err != nil) {
		fmt.Printf("Could not Read" + path)
		return
	}

	error := os.MkdirAll("public/" + outputDirectory, 0755)
	if (error != nil) {
		fmt.Printf("Could not Create Directory" + outputDirectory)
		return
	}

	out, err := os.Create("public/" + path)
	if (err != nil) {
		fmt.Printf("Could not Create File" + path)
		return
	}
	defer out.Close()

	io.Copy(out, resp.Body)
}

func fetchObject() negroni.Handler {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		defer next(rw, r)
		loadObject(r.URL.Path)
	})
}

func main() {
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.Use(cacheFor(31104000))
	n.Use(negroni.NewStatic(http.Dir("public")))
	n.Use(fetchObject())
	n.Use(negroni.NewStatic(http.Dir("public")))
	n.UseHandler(http.NotFoundHandler())
	n.Run(":8080")
}
