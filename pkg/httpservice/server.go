package httpservice

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

var muxRouter = mux.NewRouter().StrictSlash(false)
var srv *http.Server

func LaunchServer() {
	srv = &http.Server{Addr: ":3000", Handler: muxRouter}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()
	log.Println("Listening...")
}

func GET(path string, handler MyHandler) {
	router := getRouter()
	router.NewRoute().Methods("GET").Path(path).Handler(handler)
}

// waitForShutdown wait for signal from Kubelet or OS
func WaitForShutdown() {

	if srv == nil {
		return
	}

	// 1. Register a signal handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 2. Wait for kubelet to trigger shutdown
	<-sigChan
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server shutdown error: %v", err)
	}
	log.Println("HTTP server stopped.")
}

func getRouter() *mux.Router {
	return muxRouter
}
