package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// for prod serve the build from ./web/dist
	mux.Handle("/", http.FileServer(http.Dir("./web/dist")))

	log.Println("listening :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
