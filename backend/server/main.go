package main

import (
	"fmt"
	"log"
	"net/http"

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
	err = watcher.Add("./dist")
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
		rest.Get("/:context/:clientId", HTTPMeasures),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./dist"))))
	log.Fatal(http.ListenAndServe(":8090", nil))
}

// HTTPMeasures get http traffic
func HTTPMeasures(w rest.ResponseWriter, r *rest.Request) {
	ctx := r.PathParam("context")
	id := r.PathParam("clientId")
	measures := backend.AllTraffic(id, ctx)
	fmt.Printf("measures: %+v", measures)
	w.WriteJson(&measures)
}
