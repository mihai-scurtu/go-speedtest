package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Main entry point
func Run() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {

	})

	addr := ":" + os.Getenv("PORT")

	fmt.Printf("Listening on address: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
