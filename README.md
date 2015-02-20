# Image Resizer

This is a stateless, performant service, which will fetch images from a source you specify,
and resize it. This image is cached on disk for further requests.

This should ideally be fronted by a CDN like cloudfront. The syntax of the url looks like this:
http://yourserver:port/rxWIDTHxHEIGHT/RELATIVE_IMAGE_PATH/anything.EXTENSION

PS: We set a Cache-Control header for 360 days, so hopefully our load should be way low

Building:
  make

Usage:
  ./main http://golang.org/

Now check out the following urls in your browser
  http://golang.org/doc/gopher/frontpage.png
  http://localhost:8080/abcdx250x340/doc/gopher/frontpage.png/frontpage.png
  http://localhost:8080/abcdx250x340/doc/gopher/frontpage.png/frontpage.jpg
  http://localhost:8080/abcdx125x170/doc/gopher/frontpage.png/frontpage.png
  http://localhost:8080/abcdx125x0/doc/gopher/frontpage.png/frontpage.png
  http://localhost:8080/abcdx0x0/doc/gopher/frontpage.png/frontpage.png

Pending:
  Deduping requests
