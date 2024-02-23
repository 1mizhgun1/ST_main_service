package api

import (
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func StartServer(addr string) {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	corsHandler := handlers.CORS(headersOk, originsOk, methodsOk)(r)

	r.HandleFunc("/ping", handlePing)
	r.HandleFunc("/send", handleSend)
	r.HandleFunc("/test", handleTest)

	http.Handle("/", r)

	srv := http.Server{
		Handler:           corsHandler,
		Addr:              addr,
		ReadTimeout:       readTimeout * time.Second,
		WriteTimeout:      writeTimeout * time.Second,
		ReadHeaderTimeout: readHeaderTimeout * time.Second,
	}

	srv.ListenAndServe()
}
