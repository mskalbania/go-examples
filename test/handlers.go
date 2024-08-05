package test

import (
	"encoding/json"
	"net/http"
)

func Routes() {
	http.HandleFunc("/users", getUsersHandler)
}

func getUsersHandler(writer http.ResponseWriter, request *http.Request) {
	users := []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}{
		{"id1", "john"},
		{"id2", "doe"},
	}
	writer.WriteHeader(200)
	writer.Header().Add("Content-Type", "application/json")
	b, _ := json.Marshal(&users)
	_, _ = writer.Write(b)
}
