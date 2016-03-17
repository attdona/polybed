package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/attdona/polybed/db"
	"github.com/attdona/polybed/handlers"
)

var (
	dir  = "./app"
	host = "127.0.0.1"
	port = 8080
)

func init() {
	dburl := os.Getenv("POLYBED_DB_URL")
	if dburl != "" {
		db.URL = dburl
	}

	//db.GetMongoDB()
}

func main() {
	mux := http.NewServeMux()

	fileServer := http.Dir(dir)
	fileHandler := http.FileServer(fileServer)
	mux.Handle("/", fileHandler)

	mux.Handle("/user", http.HandlerFunc(handlers.UserHandler))

	log.Printf("Running on port %d\n", port)

	addr := fmt.Sprintf("%s:%d", host, port)

	err := http.ListenAndServe(addr, mux)
	fmt.Println(err.Error())
}
