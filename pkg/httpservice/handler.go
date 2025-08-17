package httpservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shamhub/exercise/pkg/errorlib"
)

type MyHandler func(*RequestContext) (interface{}, error)

func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	resourcePath := r.URL.Path

	injector := newContextInjector(r)

	injector.injectRequestContext(r)

	data, err := h(injector)
	if err != nil {
		processError(w, resourcePath, err)
		return
	}

	processData(w, data)
}

func processData(w http.ResponseWriter, data interface{}) {
	fmt.Println("processing data")
	encodedData, err := json.Marshal(data)
	if err != nil {
		marshalErrorResponse(w, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encodedData)
}

func processError(w http.ResponseWriter, resourcePath string, err error) {
	switch err := err.(type) {
	case *errorlib.ResponseError:
		c := newCustomErrorForSingleErrorResponse(resourcePath, err)
		SendErrorResponse(w, c)
	case *errorlib.MultiErrors:
		c := newCustomErrorForMultiErrorResponse(resourcePath, err)
		SendErrorResponse(w, c)
	}
}
