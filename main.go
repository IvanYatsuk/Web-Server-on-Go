package main

import (
	"fmt"
	"net/http"
)

func main_page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Main Page")
}

func main() {
	http.HandleFunc("/", main_page)
	http.ListenAndServe(":8800", nil)
}
