package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func StartServer(addr string) {
	go func() {
		err := readFromKafka()
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		ticker := time.NewTicker(kafkaReadPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				scanStorage(sendReceiveRequest)
			}
		}
	}()

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	r.Use(CORSMiddleware)

	r.HandleFunc("/send", handleSend).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/transfer", handleTransfer).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/ping", handlePing).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/test_send", handleTestSend).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/test_receive", handleTestReceive).Methods(http.MethodPost, http.MethodOptions)

	http.Handle("/", r)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	srv := http.Server{
		Handler:           r,
		Addr:              addr,
		ReadTimeout:       readTimeout * time.Second,
		WriteTimeout:      writeTimeout * time.Second,
		ReadHeaderTimeout: readHeaderTimeout * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("Server stopped")
		}
	}()
	fmt.Println("Server started")

	sig := <-signalCh
	fmt.Printf("Received signal: %v\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown failed: %v\n", err)
	}
}
