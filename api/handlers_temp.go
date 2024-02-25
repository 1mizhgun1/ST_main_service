package api

import (
	"fmt"
	"net/http"
)

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
	w.WriteHeader(http.StatusOK)
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data := codeRequest{}
	err := getRequestData(w, r, &data)
	if err != nil {
		return
	}
	
	fmt.Printf("%v\n", data)
}
