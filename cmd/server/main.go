package main

import (
	"log"
	"net"
	"net/http"

	"github.com/L0rd1k/uprise-api/cmd/api"
	"github.com/L0rd1k/uprise-api/cmd/api/router"
	"github.com/L0rd1k/uprise-api/test"
)

func TestDemo() {
	tst_api := api.NewApi()
	tst_api.Use(api.DefaultCommonStack...)

	tst_api.SetApp(
		api.AppSimple(func(w api.ResponseWriter, r *api.Request) {
			w.WriteJson(map[string]string{"Body": "Hello World!"})
		}),
	)

	log.Fatal(http.ListenAndServe(":8080", tst_api.MakeHandler()))
}

func TestDemoPlaceholder() {
	tst_api := api.NewApi()
	tst_api.Use(api.DefaultCommonStack...)
	_router, _err := router.MakeRouter(
		router.Get("/lookup/#host", func(w api.ResponseWriter, req *api.Request) {
			ip, err := net.LookupIP(req.PathParam("host"))
			if err != nil {
				api.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteJson(&ip)
		}),
	)
	if _err != nil {
		log.Fatal(_err)
	}
	tst_api.SetApp(_router)
	log.Fatal(http.ListenAndServe(":8080", tst_api.MakeHandler()))

}

func main() {
	test.TestCircleQueue()
	// test.TestArrayStack()
	// test.TestArray()
	// TestDemoPlaceholder()
	// TestDemo()
}
