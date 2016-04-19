package main

import (
	"log"
	"net/http"
	"time"

	"github.com/attdona/polybed/backend"

	"github.com/ant0ine/go-json-rest/rest"
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
	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./app"))))
	log.Fatal(http.ListenAndServe(":8090", nil))
}

// HTTPMeasures get http traffic
func HTTPMeasures(w rest.ResponseWriter, r *rest.Request) {

	filters := backend.TrafficMeasureFilter {
		Context: 	r.PathParam("context"),
		Pool: 		r.PathParam("clientId"),
		FromDate:	time.Time{},
		ToDate:		time.Time{},
	}

	measures := backend.GetGraphData(filters)
	w.WriteJson(&measures)
}
