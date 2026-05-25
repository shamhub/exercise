package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/shamhub/exercise/pkg/envconfigreader"
	"github.com/shamhub/exercise/pkg/httpservice/proxy"
	"github.com/shamhub/exercise/types"
)

func main() {

	// 1. Read config from .env files and override with shell variables
	var envReader envconfigreader.IEnvconfigReader
	envReader = envconfigreader.NewEnvConfig("./configs")

	// appName := envReader.Get(types.APP_NAME)
	kindServiceList := envReader.Get(types.SERVICE_LIST)

	serverCert := envReader.Get(types.TLS_SERVER_CERT)
	serverKey := envReader.Get(types.TLS_SERVER_KEY)

	// 3. Initialize serverpool
	serverPool := proxy.NewKindServicePool(envReader)
	healthCheckCtx, cancel := context.WithCancel(context.Background())

	targets := strings.Split(kindServiceList, ",")

	for _, target := range targets {
		targetUrl, err := url.Parse(strings.TrimSpace(target))
		if err != nil {
			panic(err)
		}

		service := proxy.NewKindService(targetUrl)
		serverPool.AddService(service)
	}

	// 4. Add healthchecker for all kindServices
	go func(ctx context.Context) {
		t := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				serverPool.HealthCheck()
			}
		}
	}(healthCheckCtx)

	// 3. Register handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", serverPool.HandleRoutes())

	// 4. Start server
	server := http.Server{
		Addr:         ":443",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	server.ListenAndServeTLS(serverCert, serverKey)

	// 5. Shutdown logic
	WaitForShutDown(cancel, &server)
}

func WaitForShutDown(cancelHealthCheck context.CancelFunc, server *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	cancelHealthCheck()

	shutdownDeadline := 30 * time.Second
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownDeadline)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("forcing shutdown - %s", err.Error())
		os.Exit(1)
	}
}
