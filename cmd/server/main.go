package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/L0rd1k/uprise-api/cmd/api"
	"github.com/L0rd1k/uprise-api/experimental/containers/lists/array"
	"github.com/L0rd1k/uprise-api/experimental/containers/lists/singlell"
	"github.com/L0rd1k/uprise-api/experimental/containers/utils"
)

func TestArray() {
	tstArray := array.New("Flone", "Apple", "Diving", "Bucket")
	fmt.Println(tstArray.ToString())
	tstArray.Sort(utils.Comparator_String)
	fmt.Println("\n", tstArray.ToString())
}

func TestSingleLinkedList() {
	tstList := singlell.New("One", "Two", "Three", "Four")
	tstList.Remove(3)
	fmt.Println(tstList.ToString())
}

func TestDemo() {
	tst_api := api.NewApi()
	tst_api.Use(api.DefaultCommonStack...)
	tst_api.SetApp(api.AppSimple(func(w api.ResponseWriter, r *api.Request) {
		w.WriteJson(map[string]string{"Body": "Hello World!"})
	}))
	log.Fatal(http.ListenAndServe(":8080", tst_api.MakeHandler()))
}

func main() {

	TestDemo()
}
