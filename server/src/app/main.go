package main

import (
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

var helloWorld = func(w http.ResponseWriter, r *http.Request) {
	log.Println("log")
	log.Println("Header", r.Header.Get("X-TOKEN"))
	fmt.Fprint(w, "Hello!")
}

func multiMiddleware(f http.HandlerFunc, m ...middleware) http.HandlerFunc {
	if len(m) < 1 {
		return f
	}

	return m[0](multiMiddleware(f, m[1:cap(m)]...))
}

func main() {
	httpPipeline := []middleware{
		logMiddleware,
		appenderMiddleware,
		authenticationMiddleware,
	}

	http.HandleFunc("/", (appenderMiddleware(authenticationMiddleware(helloWorld))))
	http.HandleFunc("/2", multiMiddleware(helloWorld))
	http.HandleFunc("/3", multiMiddleware(helloWorld, httpPipeline...))

	log.Println("Now server is running on port 5000")
	http.ListenAndServe(":5000", nil)
}
