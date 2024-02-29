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

	"github.com/1mizhgun1/ST_main_service/api/consts"
	"github.com/1mizhgun1/ST_main_service/api/handlers"
	"github.com/1mizhgun1/ST_main_service/api/kafka"
	"github.com/1mizhgun1/ST_main_service/api/middleware"
	"github.com/1mizhgun1/ST_main_service/api/storage"
	"github.com/1mizhgun1/ST_main_service/api/utils"
)

func StartServer(addr string) {
	go func() {
		err := kafka.ReadFromKafka()
		if err != nil {
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
		Addr:              addr,
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
