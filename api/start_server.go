package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func StartServer(addr string) {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	r.HandleFunc("/ping", handlePing)
	r.HandleFunc("/send", handleSend)
	r.HandleFunc("/test", handleTest)

	http.Handle("/", r)

	srv := http.Server{
		Handler:           r,
		Addr:              addr,
		ReadTimeout:       readTimeout * time.Second,
		WriteTimeout:      writeTimeout * time.Second,
		ReadHeaderTimeout: readHeaderTimeout * time.Second,
	}

	srv.ListenAndServe()
}
