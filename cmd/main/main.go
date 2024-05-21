package main

import (
	"context"
	"fmt"
	"github.com/1mizhgun1/ST_main_service/internal/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/1mizhgun1/ST_main_service/internal/consts"
	_ "github.com/1mizhgun1/ST_main_service/internal/docs"
	"github.com/1mizhgun1/ST_main_service/internal/handlers"
	"github.com/1mizhgun1/ST_main_service/internal/kafka"
	"github.com/1mizhgun1/ST_main_service/internal/storage"
	"github.com/1mizhgun1/ST_main_service/internal/utils"
)

const address = consts.MyHost + ":8080"

// @title 			Transport Layer API
// @version 		1.0
// @description 	API for Transport Layer
func main() {
	go func() {
		if err := kafka.ReadFromKafka(); err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		ticker := time.NewTicker(consts.KafkaReadPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				storage.ScanStorage(utils.SendReceiveRequest)
			}
		}
	}()

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	r.Use(middleware.CORSMiddleware)

	r.PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/send", handlers.HandleSend).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/transfer", handlers.HandleTransfer).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/ping", handlers.HandlePing).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/test_send", handlers.HandleTestSend).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/test_receive", handlers.HandleTestReceive).Methods(http.MethodPost, http.MethodOptions)

	http.Handle("/", r)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	srv := http.Server{
		Handler:           r,
		Addr:              address,
		ReadTimeout:       consts.ReadTimeout * time.Second,
		WriteTimeout:      consts.WriteTimeout * time.Second,
		ReadHeaderTimeout: consts.ReadHeaderTimeout * time.Second,
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
