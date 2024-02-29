package api

import (
	"fmt"
	"net/http"
)

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"message":"pong"}`))
	w.WriteHeader(http.StatusOK)
}

func handleTestSend(w http.ResponseWriter, r *http.Request) {
	data := receiveRequest{}
	err := getRequestData(r, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("%v\n", data)
	w.WriteHeader(http.StatusOK)
}

func handleTestReceive(w http.ResponseWriter, r *http.Request) {
	data := receiveRequest{}
	err := getRequestData(r, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("%v\n", data)
	w.WriteHeader(http.StatusOK)
}
