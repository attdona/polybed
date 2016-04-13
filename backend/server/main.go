package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/attdona/polybed/backend"
)

func main() {

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/http/:clientId", HTTPMeasures),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8090", api.MakeHandler()))
}

// HTTPMeasures get http traffic
func HTTPMeasures(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("clientId")
	fmt.Println("HITTED")
	measures := backend.AllTraffic(id, "http")
	fmt.Println("measures: ", measures)
	w.WriteJson(&measures)
}
