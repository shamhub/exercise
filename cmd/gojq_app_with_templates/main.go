package main

import (
	"github.com/shamhub/exercise/cmd/gojq_app_with_templates/handler"
)

func main() {

	// 3. Implement handlers
	handler.GET("/api/v1/user/{userId}", handler.GetUserId)
}
