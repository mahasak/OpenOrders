package middleware

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

var AuthenticationMiddleware = func(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Authentication")
		f.ServeHTTP(w, r)
	})
}

var LogMiddleware = func(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Logger")
		h.ServeHTTP(w, r)
		log.SetOutput(os.Stdout)
		log.Println(r.Method, r.Host, r.URL)
	})
}

var AppenderMiddleware = func(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("X-TOKEN", "1234")
		fmt.Println("Appender")
		h.ServeHTTP(w, r)
		log.SetOutput(os.Stdout)
		log.Println(r.Method, r.Host, r.URL)
	})
}

var EnforceXMLHandler = func(next http.Handler) http.Handler {
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
