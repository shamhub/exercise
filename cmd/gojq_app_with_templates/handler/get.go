package handler

import (
	"github.com/gorilla/mux"
	"github.com/shamhub/exercise/pkg/httpservice"
)

var muxRouter = mux.NewRouter()

func GET(path string, handler httpservice.MyHandler) {
	muxRouter.NewRoute().Methods("GET").Path(path).Handler(handler)
}

func GetUserId(ctx *httpservice.RequestContext) (interface{}, error) {

	return nil, nil
}
