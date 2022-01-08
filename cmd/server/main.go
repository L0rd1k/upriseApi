package main

import (
	"fmt"
	"net/http"
)

type CustomTestStruct struct {
	name float64
	arr  []int
}

func main() {
	// http.HandleFunc("/", customHandler)
	// log.Fatal(http.ListenAndServe(":8080", nil))
}

func customHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "My Title", "My Body")
}
