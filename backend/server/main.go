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
		rest.Get("/:context/:clientId", HTTPMeasures),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./app"))))
	log.Fatal(http.ListenAndServe(":8090", nil))
}

// HTTPMeasures get http traffic
func HTTPMeasures(w rest.ResponseWriter, r *rest.Request) {
	ctx := r.PathParam("context")
	id := r.PathParam("clientId")
	fmt.Println("HITTED")
	measures := backend.AllTraffic(id, ctx)
	fmt.Println("measures: ", measures)
	w.WriteJson(&measures)
}