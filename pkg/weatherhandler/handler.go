package handler

import (
	"strconv"

	"github.com/gorilla/schema"
	"github.com/shamhub/exercise/pkg/errorlib"
	"github.com/shamhub/exercise/pkg/httpservice"
	"github.com/shamhub/exercise/pkg/weatherservice"
)

type DTO struct {
	Weather     string `json:"weather"`
	Temperature string `json:"temperature"`
}

type queryParamStruct struct {
	Latitude  string `schema:"lat"`
	Longitude string `schema:"long"`
}

func GetWeatherDetails(ctx *httpservice.RequestContext) (dto interface{}, err error) {

	dto = []DTO{}

	// 1. Read query params
	var queryParams map[string][]string
	queryParams, err = ctx.GetQueryParams()
	if err != nil {
		err = errorlib.NewResponseError(400, err.Error())
		return
	}

	// 2. Validate query params
	var decoder = schema.NewDecoder()
	var params queryParamStruct
	err = decoder.Decode(&params, queryParams)
	if err != nil {
		err = errorlib.NewResponseError(400, err.Error())
		return
	}
	err = errorlib.ValidateQueryParams(params)
	if err != nil {
		return
	}
	_, err = strconv.ParseFloat(params.Latitude, 64)
	if err != nil {
		err = errorlib.NewResponseError(400, err.Error())
		return
	}
	_, err = strconv.ParseFloat(params.Longitude, 64)
	if err != nil {
		err = errorlib.NewResponseError(400, err.Error())
		return
	}

	// 3. process latitude,longitude inputs and provide weather data
	hostPortURL := "https://api.weather.gov"
	wservice := weatherservice.NewWeatherService(params.Latitude, params.Longitude, hostPortURL)

	var weatherDetails weatherservice.WeatherDetails
	pointsApi := "points/" + params.Latitude + "," + params.Longitude
	weatherDetails, err = wservice.GetWeatherDetail(pointsApi)
	if err != nil {
		return
	}

	temperatureInCelsius := strconv.FormatFloat(weatherDetails.Temperature, 'f', -1, 64) + "°C"
	dto = []DTO{
		{
			Temperature: temperatureInCelsius,
			Weather:     weatherDetails.Weather,
		},
	}
	return
}
