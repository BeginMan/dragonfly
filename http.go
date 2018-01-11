package main

import (
	"net/http"
	"encoding/json"
)

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		res := Res{"", 404, "404 not found"}
		jsonRes, _ := json.Marshal(res)
		w.Write(jsonRes)
	}
}