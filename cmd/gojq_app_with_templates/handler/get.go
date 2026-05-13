package handler

import "github.com/shamhub/exercise/pkg/httpservice"

func GetUserId(ctx *httpservice.RequestContextForTemplate) (interface{}, error) {
	return "sham", nil
}
