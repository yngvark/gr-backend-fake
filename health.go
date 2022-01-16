package main

import (
	"fmt"
	"net/http"
)

func health(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "OK")
}
