package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))

	serv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := serv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
