package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"sync/atomic"
	"testing"

	"github.com/shamhub/exercise/pkg/httpservice/proxy"
)

func setupTestService(target string, alive bool, connections int64) *proxy.KindService {
	URL, _ := url.Parse(target)
	return &proxy.KindService{
		URL:         URL,
		Alive:       alive,
		ActiveConns: connections,
		ReverseProxy: &httputil.ReverseProxy{
			Director: func(req *http.Request) {},
		},
	}
}

// 1. verify least connection behaviour
func TestGetNextBackend(t *testing.T) {
	b1 := setupTestService("http://server_busy", true, 25)
	b2 := setupTestService("http://server_free", true, 1)
	b3 := setupTestService("http://server_offline", false, 0)

	servicePool := proxy.KindServicePool{
		KindServiceList: []*proxy.KindService{b1, b2, b3},
		Log:             log.New(os.Stdout, "test", log.Ldate|log.Ltime|log.Lshortfile),
	}

	gotService := servicePool.GetNextAvailableService()
	if gotService == nil || gotService.URL.String() != "http://server_free" {
		t.Errorf("Expected lowest work server 'http://server_free' got: %v", gotService)
	}
}

func TestHandleRoutes_SuccessRoutingAndMetrics(t *testing.T) {
	// Initialize target test metadata
	targetURL, _ := url.Parse("http://mock-internal-service:8080")
	mockService := &proxy.KindService{
		URL:   targetURL,
		Alive: true, // Must be true for GetNextAvailableService loop to pick it
	}

	// Capture values inside the proxy runtime to verify atomic state transitions
	var proxyExecuted bool
	var inFlightConns int64
	var capturedPath string

	mockService.ReverseProxy = &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			proxyExecuted = true
			capturedPath = req.URL.Path
			// Snapshot active connections DURING the middle of execution
			inFlightConns = atomic.LoadInt64(&mockService.ActiveConns)
		},
	}

	// Setup the pool with our test service using a discarded logger to suppress test noise
	ksPool := &proxy.KindServicePool{
		Log: log.New(io.Discard, "", 0),
	}
	ksPool.AddService(mockService)

	// Execute mock request using httptest standard modules
	req := httptest.NewRequest("GET", "/api/v1/orders?id=99", nil)
	rec := httptest.NewRecorder()

	handler := ksPool.HandleRoutes()
	handler.ServeHTTP(rec, req)

	// Assert: Check that proxy wrapper was successfully invoked
	if !proxyExecuted {
		t.Fatal("Expected ReverseProxy to be executed, but the route step was skipped")
	}

	// Assert: Path passed cleanly down to backend destination intact
	if capturedPath != "/api/v1/orders" {
		t.Errorf("Expected request path to flow untouched, got %q", capturedPath)
	}

	// Assert: Counter increased to exactly 1 inside the execution window frame
	if inFlightConns != 1 {
		t.Errorf("Expected ActiveConns tracking variable to hit 1 during processing, found: %d", inFlightConns)
	}

	// Assert: Post-execution logic executed defer loop cleanly, resetting work pool balancing state
	if atomic.LoadInt64(&mockService.ActiveConns) != 0 {
		t.Errorf("Expected ActiveConns metric cleanup to decrement back to 0, found: %d", mockService.ActiveConns)
	}
}
