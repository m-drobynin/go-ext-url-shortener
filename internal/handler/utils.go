package handler

import "net/http"

func handleNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))
}

func handleBadRequest(w http.ResponseWriter, err *error) {
	w.WriteHeader(http.StatusBadRequest)

	if err != nil {
		var error = *err
		w.Write([]byte(error.Error()))
	} else {
		w.Write([]byte("bad request"))
	}
}
