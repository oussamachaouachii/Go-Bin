package main

import (
	"log"
	"net/http"
)

func main() {
	//addr := os.Getenv("SNIPPETBOX_ADDR")
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Println("Starting server on http://localhost:9000")
	err := http.ListenAndServe(":9000", mux)
	log.Fatal(err)
}
