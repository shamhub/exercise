package main

import (
	"github.com/shamhub/exercise/pkg/httpservice"
	handler "github.com/shamhub/exercise/pkg/weatherhandler"
)

func main() {
	// 1. Map api url's
	httpservice.GET("/weather", handler.GetWeatherDetails)

	// 2. start http server
	httpservice.LaunchServer()

	// 3. handle shutdown logic
	httpservice.WaitForShutdown()
}
