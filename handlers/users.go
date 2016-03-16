package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func UserHandler(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "POST":
		response.Header().Set("Content-type", "application/json")
		// json data to send to client

		decoder := json.NewDecoder(request.Body)
		t := make(map[string]string)
		err := decoder.Decode(&t)
		if err != nil {
			fmt.Errorf("Oops! Error: %s\n", err)
		}

		fmt.Printf("Request Body: %+v\n", t)

		name, ok := t["name"]
		if !ok {
			name = "me, myself and I"
		}

		api, aok := t["api"]
		if !aok {
			api = "user"
		}

		data := map[string]string{"api": api, "name": name}

		json_bytes, _ := json.Marshal(data)
		fmt.Printf("json_bytes: %s\n", string(json_bytes[:]))
		fmt.Fprintf(response, "%s\n", json_bytes)
	default:
		http.Error(response, "Invalid request method.", 405)
	}
}
