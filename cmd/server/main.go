package main

import (
	"fmt"
	"net/http"

	"github.com/L0rd1k/uprise-api/experimental/containers/lists/array"
	"github.com/L0rd1k/uprise-api/experimental/containers/utils"
)

func main() {
	tstArray := array.New("Flone", "Apple", "Diving", "Bucket")
	fmt.Println(tstArray.ToString())
	tstArray.Sort(utils.Comparator_String)

	fmt.Println("\n", tstArray.ToString())
	// fmt.Println(myList)
	// http.HandleFunc("/", customHandler)
	// log.Fatal(http.ListenAndServe(":8080", nil))
}

func customHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "My Title", "My Body")
}
