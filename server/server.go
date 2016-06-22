package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Main entry point
func Run() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		// Number of bytes to download
		size, err := strconv.Atoi(r.FormValue("size"))

		if err != nil || size < 1 || size > 1000000000 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Send buffer
		buf := make([]byte, size)
		w.Write(buf)
	})

	addr := ":" + os.Getenv("PORT")

	fmt.Printf("Listening on address: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
