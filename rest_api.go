package main

import (
	"fmt"
	"log"
	"net/http"
)

func root(w http.ResponseWriter, req *http.Request) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	fmt.Println("Request to /")
	_, _ = fmt.Fprintf(w, "Welcome to the home page!")
}

func health(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "OK")
}

func handleRequests(port int) {
	var address string = fmt.Sprintf("%v%v", ":", port)

	http.HandleFunc("/", root)
	http.HandleFunc("/health", health)
	log.Fatal(http.ListenAndServe(address, nil))
}

func main() {
	port := 8080
	fmt.Printf("Running web server on port %v\n", port)

	handleRequests(port)
}
