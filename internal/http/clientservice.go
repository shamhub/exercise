package http

import (
	"context"
	"net/http"
	"time"
)

// http client per target host
//   - currently handles JSON,xml,plain/text responses
//   - currently support request with token based Authorization headers
type HttpClientService struct {
	*http.Client
	hostPortURL          string
	serviceLevelOptions  *ServiceLevelOptions
	contentTypeSupported responseType
}

type responseType int

const (
	JSON responseType = iota
	XML
	TEXT
	UNSUPPORTED_RESPONSE_TYPE
)

func NewHTTPServiceWithOptions(hostportURL string) *HttpClientService {

	createHttpClient := func() *http.Client {

		// custom HTTP client with the logging middleware
		return &http.Client{
			Transport: NewLoggingRoundTripper(),
			// The total time allowed for an HTTP client to complete a request,
			// including connection establishment, sending the request, and receiving
			// the response. This timeout encompasses all phases of a single client-side operation.
			Timeout: time.Duration(60) * time.Second,
		}
	}

	return &HttpClientService{
		hostPortURL: hostportURL,
		// custom HTTP client with the logging middleware
		Client:              createHttpClient(),
		serviceLevelOptions: NewServiceLevelOption(),
	}
}

func (h *HttpClientService) Get(ctx context.Context, api string, queryParams map[string]interface{}) (*Response, error) {
	return h.GetWithHeaders(ctx, api, queryParams, map[string]string{})
}
func (h *HttpClientService) GetWithHeaders(ctx context.Context, api string, queryParams map[string]interface{},
	headers map[string]string) (*Response, error) {
	return h.callService(ctx, http.MethodGet, api, queryParams, nil, headers)
}
