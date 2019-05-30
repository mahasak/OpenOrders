package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("log")
		fmt.Fprint(w, "Hello!")
	})

	log.Println("Now server is running on port 3000")
	http.ListenAndServe(":3000", nil)
}
