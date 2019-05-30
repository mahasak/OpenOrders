package main

import (
	"fmt"
	"log"
	"middleware"
	"net/http"
)

var finalHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
})

var helloWorld = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	log.Println("log")
	log.Println("Header", r.Header.Get("X-TOKEN"))
	fmt.Fprint(w, "Hello!")
})

func main() {
	privatePipeline := []middleware.Handler{middleware.LogMiddleware, middleware.AppenderMiddleware, middleware.AuthenticationMiddleware}
	publicPipeline := []middleware.Handler{middleware.LogMiddleware, middleware.AppenderMiddleware}

	http.Handle("/public", middleware.Chain(publicPipeline...).Then(finalHandler))
	http.Handle("/private", middleware.Chain(privatePipeline...).Then(finalHandler))

	log.Println("Now server is running on port 5000")
	http.ListenAndServe(":5000", nil)
}
