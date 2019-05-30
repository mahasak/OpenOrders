package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

type middleware func(http.HandlerFunc) http.HandlerFunc

var authenticationMiddleware = func(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Authentication")
		f.ServeHTTP(w, r)
	}
}

var logMiddleware = func(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Logger")
		h.ServeHTTP(w, r)
		log.SetOutput(os.Stdout)
		log.Println(r.Method, r.Host, r.URL)
	}
}

var appenderMiddleware = func(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("X-TOKEN", "1234")
		fmt.Println("Appender")
		h.ServeHTTP(w, r)
		log.SetOutput(os.Stdout)
		log.Println(r.Method, r.Host, r.URL)
	}
}

var enforceXMLHandler = func(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for a request body
		if r.ContentLength == 0 {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		// Check its MIME type
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		if http.DetectContentType(buf.Bytes()) != "text/xml; charset=utf-8" {
			http.Error(w, http.StatusText(415), 415)
			return
		}
		next.ServeHTTP(w, r)
	})
}

var finalHandler = func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

var helloWorld = func(w http.ResponseWriter, r *http.Request) {
	log.Println("log")
	log.Println("Header", r.Header.Get("X-TOKEN"))
	fmt.Fprint(w, "Hello!")
}

func multiMiddleware(f http.HandlerFunc, m ...middleware) http.HandlerFunc {
	if len(m) == 0 {
		return f
	}

	return m[0](multiMiddleware(f, m[1:cap(m)]...))
}

func main() {
	privateHTTPPipeline := []middleware{
		logMiddleware,
		appenderMiddleware,
		authenticationMiddleware,
	}

	publicHTTPPipeline := []middleware{
		logMiddleware,
		appenderMiddleware,
	}

	http.HandleFunc("/", logMiddleware(appenderMiddleware(authenticationMiddleware(helloWorld))))
	http.HandleFunc("/public", multiMiddleware(helloWorld, publicHTTPPipeline...))
	http.HandleFunc("/private", multiMiddleware(helloWorld, privateHTTPPipeline...))
	http.HandleFunc("/xml", enforceXMLHandler(finalHandler))

	log.Println("Now server is running on port 5000")
	http.ListenAndServe(":5000", nil)
}
