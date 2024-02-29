package handlers

import (
	"fmt"
	"net/http"

	"github.com/1mizhgun1/ST_main_service/api/utils"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"message":"pong"}`))
	w.WriteHeader(http.StatusOK)
}

func HandleTestSend(w http.ResponseWriter, r *http.Request) {
	data := utils.ReceiveRequest{}
	err := utils.GetRequestData(r, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("%v\n", data)
	w.WriteHeader(http.StatusOK)
}

func HandleTestReceive(w http.ResponseWriter, r *http.Request) {
	data := utils.ReceiveRequest{}
	err := utils.GetRequestData(r, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("%v\n", data)
	w.WriteHeader(http.StatusOK)
}
