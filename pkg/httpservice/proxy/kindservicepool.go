package proxy

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/shamhub/exercise/pkg/envconfigreader"
	"github.com/shamhub/exercise/types"
)

// Maintains the list of KindServices
type KindServicePool struct {
	KindServiceList []*KindService
	Log             *log.Logger
}

func NewKindServicePool(config envconfigreader.IEnvconfigReader) *KindServicePool {
	return &KindServicePool{
		Log: log.New(os.Stdout, config.Get(types.APP_NAME), log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (ksPool *KindServicePool) AddService(b *KindService) {
	ksPool.KindServiceList = append(ksPool.KindServiceList, b)
}

// GetNextAvailableService() checks whether service is alive with minimum connections among others
func (ksPool *KindServicePool) GetNextAvailableService() *KindService {
	var minimumConnections int64 = -1
	var chosenKindService *KindService
	for _, kindService := range ksPool.KindServiceList {

		if !kindService.IsAlive() {
			continue
		}

		activeConnections := atomic.LoadInt64(&kindService.ActiveConns)

		if minimumConnections == -1 || activeConnections < minimumConnections {
			minimumConnections = activeConnections
			chosenKindService = kindService
		}
	}
	return chosenKindService
}

// HealthCheckError need to be logged
func (ksPool *KindServicePool) HealthCheck() {
	for _, kindService := range ksPool.KindServiceList {
		alive, err := IsKindServiceAlive(kindService.URL)
		if err != nil {
			msg := fmt.Sprintf("HealthCheck() error - %s", err.Error())
			log.Println(msg)
		}

		kindService.SetAliveness(alive)
	}
}

func (ksPool *KindServicePool) HandleRoutes() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// 1. get the active kindService
		kindService := ksPool.GetNextAvailableService()
		if kindService == nil {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}

		atomic.AddInt64(&kindService.ActiveConns, 1)
		defer atomic.AddInt64(&kindService.ActiveConns, -1)

		// 2. Forward the request
		kindService.ReverseProxy.ServeHTTP(w, r)
	}
}
