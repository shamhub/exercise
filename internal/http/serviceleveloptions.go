package http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

// service level options for http client per target host
type ServiceLevelOptions struct {
	ServiceLevelHeaders map[string]string

	// numOfRetries int
	// cb *CircuitBreaker
	// *OAuthOption
	// *SurgeProtectorOption
}

func NewServiceLevelOption() *ServiceLevelOptions {
	return &ServiceLevelOptions{
		ServiceLevelHeaders: map[string]string{
			"Accept": "application/json,application/xml,text/plain",
		},
	}
}

func (s *ServiceLevelOptions) setServiceLevelHeaders(req *http.Request, body []byte) {
	if body == nil {
		return
	}

	contentType := "text/plain"
	var t interface{}

	err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&t)
	if err == nil {
		contentType = "application/json"
	}

	err = xml.NewDecoder(bytes.NewBuffer(body)).Decode(&t)
	if err == nil {
		contentType = "application/xml"
	}
	req.Header.Add("content-type", contentType)
}
func setTransport() *http.Transport {
	maxIdleConnections := 100

	maxIdleConnectionsPerhost := 100

	httpIdleConnectionTimeout := 30

	return &http.Transport{
		// MaxIdleConns is the connection pool size
		// Maximum number of idle (keep-alive) connections across all hosts

		// Increasing these values allows your HTTP client to reuse existing connections
		// for subsequent requests to the same host, which can improve performance by
		// reducing the overhead of establishing new connections for each request.
		MaxIdleConns: maxIdleConnections,

		// Maximum number of idle connections per host
		// MaxIdleConnsPerHost is helpful when you're dealing with a high
		// number of requests to a small number of hosts and want to ensure
		// efficient reuse of connections for those hosts.

		// For example, if you have many concurrent requests to the same host,
		// increasing MaxIdleConnsPerHost can reduce the overhead of
		// establishing new connections according to the official documentation.

		// Increasing these values allows your HTTP client to reuse existing connections
		// for subsequent requests to the same host, which can improve performance by
		// reducing the overhead of establishing new connections for each request.
		MaxIdleConnsPerHost: maxIdleConnectionsPerhost,

		// IdleConnTimeout is the maximum amount of time an idle (keep-alive) connection
		// will remain idle before closing itself.
		IdleConnTimeout: time.Duration(httpIdleConnectionTimeout) * time.Second,

		// DisableKeepAlives: true setting in the http.Transport prevents the client from
		// reusing the connection for subsequent requests, forcing a new connection to be
		// established for each request.
		DisableKeepAlives: false,
	}
}

// loggingRoundTripper implements http.RoundTripper for logging requests/responses.
type loggingRoundTripper struct {
	next http.RoundTripper
	log  *log.Logger
}

func NewLoggingRoundTripper() *loggingRoundTripper {

	return &loggingRoundTripper{
		next: setTransport(),
		log:  log.New(os.Stdout, "weatherservice:", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// RoundTrip implements the http.RoundTripper interface.
// The RoundTrip() method takes an *http.Request as input and is responsible for
// sending this request, receiving the response, and returning an *http.Response and an error.
// By implementing the RoundTripper interface, developers can customize the behavior of http.Client instances. This allows for the creation of "middleware" for HTTP clients, enabling features like:
// 1. Adding custom headers
// 2. Logging requests and responses
// 3. Implementing retry mechanisms
// 4. Handling authentication
// 5. Mocking HTTP responses for testing
// Note: net/http internally invokes rt.RoundTrip(req) for every client.Do(req)
func (lrt *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var requestMethod string
	var reqURL string
	if req != nil {
		requestMethod = req.Method
		reqURL = req.URL.String()
		reqURL = strings.ReplaceAll(reqURL, "\n", "")
		reqURL = strings.ReplaceAll(reqURL, "\r", "")
	}
	log.Println("HTTPRequestDetails", zap.String("reqMethod", requestMethod),
		zap.String("RequestURL", reqURL))

	start := time.Now()

	resp, err := lrt.next.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start).Seconds()
	log.Println("HTTPResponseDetails", zap.String("responseStatus", resp.Status), zap.Float64("duration", duration))
	return resp, nil
}
