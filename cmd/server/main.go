package main

import (
	"fmt"

	"github.com/L0rd1k/uprise-api/cmd/api"
)

func main() {

	/*** ARRAY
		tstArray := array.New("Flone", "Apple", "Diving", "Bucket")
		fmt.Println(tstArray.ToString())
		tstArray.Sort(utils.Comparator_String)
		fmt.Println("\n", tstArray.ToString())
	***/

	/*** SINGLE-LINKED-LIST
		tstList := singlelinkedlist.New("One", "Two", "Three", "Four")
		tstList.Remove(3)
		fmt.Println(tstList.ToString())
		fmt.Println(myList)
	***/
	tst_api := api.NewApi()
	fmt.Println(tst_api)
	fmt.Printf("Check function!")
}
