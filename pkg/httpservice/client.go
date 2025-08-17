package httpservice

import (
	"bytes"
	"context"
	"net/http"

	httpclient "github.com/shamhub/exercise/internal/http"
	httpresponse "github.com/shamhub/exercise/internal/http"
	"github.com/shamhub/exercise/pkg/errorlib"
)

type IHttpClient interface {
	Get(ctx context.Context, requestDetail *HttpRequestDetail) (*ResponseData, error)
}

type IResponse interface {
	// get response header key's value
	GetHeader(key string) string

	// get response Headers
	Headers() http.Header

	// Get response body
	GetBody() []byte

	// Get response status
	GetStatus() int

	// Bind() takes response and binds it to i based on content-type
	Bind(response []byte, i interface{}) error
}

type HttpRequestDetail struct {
	Api                 string
	QueryParams         map[string]interface{}
	Body                []byte
	RequestLevelHeaders map[string]string
}

type httpClient struct {
	*httpclient.HttpClientService
}

func NewHTTPClient(hostportURL string) *httpClient {

	return &httpClient{
		HttpClientService: httpclient.NewHTTPServiceWithOptions(hostportURL),
	}
}

func (h *httpClient) Get(ctx context.Context, requestDetail *HttpRequestDetail) (*ResponseData, error) {
	resp, err := h.HttpClientService.Get(ctx, requestDetail.Api, requestDetail.QueryParams)
	return &ResponseData{
		response: resp,
	}, err
}

type ResponseData struct {
	response *httpresponse.Response
}

func (responseData *ResponseData) GetBody() (*bytes.Reader, error) {
	if responseData == nil {
		return nil, errorlib.NewResponseError(204, "No content found")
	}
	return responseData.response.GetBody(), nil
}
