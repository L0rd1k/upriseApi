package main

import (
	"fmt"
	"net/http"

	"github.com/L0rd1k/uprise-api/experimental/containers/lists/array"
)

type CustomTestStruct struct {
	name float64
	arr  []int
}

func main() {
	myList := array.New("One", "Two", "Three")
	myList.Insert(4, "1", "2")
	fmt.Println(myList.ToString())
	// fmt.Println(myList)
	// http.HandleFunc("/", customHandler)
	// log.Fatal(http.ListenAndServe(":8080", nil))
}

func customHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "My Title", "My Body")
}
