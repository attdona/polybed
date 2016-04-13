package main

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/attdona/polybed/backend"
	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault(backend.ROPCollection, "traffic")
	viper.SetDefault(backend.DBName, "netdata")

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/:context/:clientId", HTTPMeasures),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8090", api.MakeHandler()))
}

// HTTPMeasures get http traffic
func HTTPMeasures(w rest.ResponseWriter, r *rest.Request) {
	ctx := r.PathParam("context")
	id := r.PathParam("clientId")
	measures := backend.AllTraffic(id, ctx)
	//fmt.Println("measures: ", measures)
	w.WriteJson(&measures)
}
