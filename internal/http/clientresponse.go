package http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
)

type Response struct {
	Body       *bytes.Reader
	StatusCode int
	headers    http.Header
}

func (r *Response) GetHeader(key string) string {
	if r != nil && r.headers != nil {
		return r.headers.Get(key)
	}
	return ""
}

func (r *Response) Headers() http.Header {
	return r.headers.Clone()
}

func (r *Response) GetBody() *bytes.Reader {
	if r != nil {
		return r.Body
	}
	return nil
}

func (r *Response) GetStatus() int {
	if r == nil {
		return -1
	}
	return r.StatusCode
}

// Bind() takes response and binds it to i based on content-type
func (r *Response) Bind(response []byte, i interface{}) error {

	var err error

	switch getResponseContentType(r.headers) {
	case XML:
		err = xml.NewDecoder(bytes.NewBuffer(response)).Decode(&i)
	case TEXT:
		v, ok := i.(*string)
		if ok {
			*v = fmt.Sprintf("%s", response)
		}
	case JSON:
		err = json.NewDecoder(bytes.NewBuffer(response)).Decode(&i)
	default:
		err = errors.New("unsupported response type")
	}
	return err
}

func getResponseContentType(header http.Header) responseType {
	switch header.Get("content-type") {
	case "application/xml":
		return XML
	case "text/plain":
		return TEXT
	case "application/json", "application/geo+json":
		return JSON
	default:
		return UNSUPPORTED_RESPONSE_TYPE
	}
}
