package main

import (
	"github.com/shamhub/exercise/cmd/gojq_app_with_templates/handler"
	"github.com/shamhub/exercise/pkg/httpservice"
)

func main() {

	// 1. collect rules from config
	ruleConfig := httpservice.GetActiveRules()
	if ruleConfig == nil {
		panic("rule config file is missing")
	}

	// 2. Register http handlers
	httpservice.GETTemplate("/a/b", handler.GetUserId)
	httpservice.POSTTemplate("/api/v1/user/{userId}", handler.PostUserId)

	// 3. start http server
	httpservice.LaunchServer()

	// 4. handle shutdown logic
	httpservice.WaitForShutdown()
}
