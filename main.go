package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	
	server := &http.Server {
		Handler: mux,
		Addr: ":" + port,
	}


	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}