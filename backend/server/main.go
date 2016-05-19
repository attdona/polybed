package main

import (
	"log"
	"net/http"
	"time"

	"gopkg.in/fsnotify.v1"

	"github.com/attdona/polybed/backend"
	"github.com/jaschaephraim/lrserver"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault(backend.ROPCollection, "traffic")
	viper.SetDefault(backend.DBName, "netdata")

	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	// Add dir to watcher
	err = watcher.Add("./app")
	if err != nil {
		log.Fatalln(err)
	}

	// Create and start LiveReload server
	lr, _ := lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	go lr.ListenAndServe()

	// // Start goroutine that requests reload upon watcher event
	// go func() {
	// 	for {
	// 		select {
	// 		case event := <-watcher.Events:
	// 			lr.Reload(event.Name)
	// 		case errs := <-watcher.Errors:
	// 			log.Println(errs)
	// 		}
	// 	}
	// }()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/google/:context/:clientId", HTTPMeasures),
		rest.Get("/plotly/:context/:clientId", PlotlyHTTPMeasures),
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

	filters := backend.TrafficMeasureFilter{
		Context:  r.PathParam("context"),
		Pool:     r.PathParam("clientId"),
		FromDate: time.Time{},
		ToDate:   time.Time{},
	}

	measures := backend.GetGraphData(filters)
	w.WriteJson(&measures)
}

// HTTPMeasures get http traffic
func PlotlyHTTPMeasures(w rest.ResponseWriter, r *rest.Request) {

	filters := backend.TrafficMeasureFilter{
		Context:  r.PathParam("context"),
		Pool:     r.PathParam("clientId"),
		FromDate: time.Time{},
		ToDate:   time.Time{},
	}

	measures := backend.GetPlotlyGraphData(filters)
	w.WriteJson(&measures)
}
