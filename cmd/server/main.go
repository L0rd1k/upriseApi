package main

import (
	"fmt"
	"net/http"

	singlelinkedlist "github.com/L0rd1k/uprise-api/experimental/containers/lists/single_linked_list"
)

func main() {

	/*** ARRAY
		tstArray := array.New("Flone", "Apple", "Diving", "Bucket")
		fmt.Println(tstArray.ToString())
		tstArray.Sort(utils.Comparator_String)
		fmt.Println("\n", tstArray.ToString())
	***/

	/*** SINGLE-LINKED-LIST
	***/
	tstList := singlelinkedlist.New("One", "Two", "Three", "Four")
	tstList.Remove(3)
	fmt.Println(tstList.ToString())
	// fmt.Println(myList)
	// http.HandleFunc("/", customHandler)
	// log.Fatal(http.ListenAndServe(":8080", nil))
}

func customHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "My Title", "My Body")
}
