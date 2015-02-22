package image_proxy

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/nfnt/resize"
	"github.com/golang/groupcache/singleflight"
	"net/http"
	"strings"
	"os"
	"io"
	"image"
	"image/png"
	"image/jpeg"
	"strconv"
)

func readAndResize(properties string, inputStream io.Reader) image.Image {
	img, _, err := image.Decode(inputStream)
	if(err != nil) {
		return nil
	}

	parts := strings.Split(properties, "x")

	if (len(parts) < 3) {
		return nil
	}


	width, err := strconv.ParseUint(parts[1], 10, 32)
	if (err != nil) {
		return nil
	}

	height, err := strconv.ParseUint(parts[2], 10, 32)
	if (err != nil) {
		return nil
	}

	return resize.Resize(uint(width), uint(height), img, resize.MitchellNetravali)
}


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

func loadObject(imageSource string, path string) {
	parts := strings.Split(path, "/")

	if (len(parts) < 4) {
		return
	}
	properties, key, filename := parts[1], strings.Join(parts[2:len(parts) - 1], "/"), parts[len(parts) - 1]
	outputDirectory := strings.Join(parts[1:len(parts) - 1], "/")
	extension := getExtension(filename)

	resp, err := http.Get(imageSource + key)
	defer resp.Body.Close()
	if (err != nil || resp.StatusCode != 200) {
		fmt.Printf("[ERROR] Could not Read " + path + "\n")
		return
	}

	error := os.MkdirAll("cached-images/" + outputDirectory, 0755)
	if (error != nil) {
		fmt.Printf("[ERROR] Could not Create Directory " + outputDirectory + "\n")
		return
	}

	img := readAndResize(properties, resp.Body)

	if (img == nil) {
		return
	}

	out, err := os.Create("cached-images/" + path)
	if (err != nil) {
		fmt.Printf("[ERROR] Could not Create File " + path + "\n")
		return
	}
	defer out.Close()

	switch strings.ToLower(extension) {
	case "png": png.Encode(out, img)
	case "jpg": jpeg.Encode(out, img, nil)
	case "jpeg": jpeg.Encode(out, img, nil)
	}
}

func fetchObject(imageSource string) negroni.Handler {
	group := singleflight.Group{}

	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		defer next(rw, r)
		group.Do(r.URL.Path, func () (interface{}, error) {
			loadObject(imageSource, r.URL.Path)
			return nil, nil
		})

	})
}

func main() {
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.Use(cacheFor(31104000))
	n.Use(negroni.NewStatic(http.Dir("cached-images")))
	n.Use(fetchObject(os.Args[1]))
	n.Use(negroni.NewStatic(http.Dir("cached-images")))
	n.UseHandler(http.NotFoundHandler())
	n.Run(":8080")
}
