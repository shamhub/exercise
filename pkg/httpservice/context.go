package httpservice

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type queryParamsKeyType string
type routeVarKeyType string
type payloadKeyType string

type RequestContext struct {
	ctx            context.Context
	queryParamsKey queryParamsKeyType
	routeVarKey    routeVarKeyType
	payloadKey     payloadKeyType
}

func newContextInjector(r *http.Request) *RequestContext {
	return &RequestContext{
		ctx:            context.Background(),
		queryParamsKey: "queryParams",
		routeVarKey:    "routeVariables",
		payloadKey:     "payloadKey",
	}
}

func (c *RequestContext) injectRequestContext(r *http.Request) {
	if c == nil {
		panic("initialize request context")
	}
	c.injectQueryParams(r)
	c.injectRouteVar(r)
	c.injectRequestPayLoad(r)
}

func (c *RequestContext) injectRouteVar(r *http.Request) {
	if c == nil {
		panic("initialize request context")
	}
	routeVariables := mux.Vars(r)
	c.ctx = context.WithValue(c.ctx, c.routeVarKey, routeVariables)
}

func (c *RequestContext) GetRouteVars() (map[string]string, error) {
	if c == nil {
		panic("initialize request context")
	}
	val, ok := c.ctx.Value(c.routeVarKey).(map[string]string)
	if !ok {
		return map[string]string{}, errors.New("route vars not found in context")
	}
	return val, nil
}

func (c *RequestContext) injectRequestPayLoad(r *http.Request) {
	if c == nil {
		panic("initialize request context")
	}

	injectBodyIntoContext := func(body io.ReadCloser) {
		c.ctx = context.WithValue(c.ctx, c.payloadKey, body)
	}

	// You always need to read the body to know what the contents are.
	// The client could send the body in chunked encoding with no Content-Length,
	// or it could even have an error and send a Content-Length and no body.
	// The client is never obligated to send what it says it's going to send.

	// The EOF check can work if you're only checking for the empty body, but
	// I would still also check for other error cases besides the EOF string.
	injectBodyIntoContext(r.Body)
}

func (c *RequestContext) GetRequestPayload() io.ReadCloser {
	if c == nil {
		panic("initialize request context")
	}
	val := c.ctx.Value(c.payloadKey).(io.ReadCloser)
	return val
}

func (c *RequestContext) injectQueryParams(r *http.Request) {
	if c == nil {
		panic("initialize request context")
	}
	queryParams := r.URL.Query()
	c.ctx = context.WithValue(c.ctx, c.queryParamsKey, queryParams)
}

func (c *RequestContext) GetQueryParams() (map[string][]string, error) {
	if c == nil {
		panic("initialize request context")
	}
	val, ok := c.ctx.Value(c.queryParamsKey).(url.Values)
	if !ok {
		return map[string][]string{}, errors.New("query parameters not found in context.")
	}

	return val, nil
}
