package weatherservice

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/shamhub/exercise/pkg/errorlib"
	"github.com/shamhub/exercise/pkg/httpservice"
)

type WeatherService struct {
	Latitude  string
	Longitude string
	client    httpservice.IHttpClient
}

func NewWeatherService(latitude, longitude, hostPortURL string) *WeatherService {
	return &WeatherService{
		Latitude:  latitude,
		Longitude: longitude,
		client:    httpservice.NewHTTPClient(hostPortURL),
	}
}

func (ws *WeatherService) GetWeatherDetail(pointsApi string) (weatherDetails WeatherDetails, err error) {

	var endPoint string

	// 1. get grid end point with coordinates
	endPoint, err = ws.getGridForecastEndPoint(pointsApi)
	if err != nil {
		return
	}

	var gridCoordinates string
	gridCoordinates, err = findGridCoordinates(endPoint)
	if err != nil {
		return
	}
	forecastDataApi := "gridpoints/TOP/" + gridCoordinates

	// 2. get forecast for grid coordinates
	var gridData gridForeCastData
	gridData, err = ws.getGridForecastData(forecastDataApi)
	if err != nil {
		return
	}

	// 3.
	// > Returns the short forecast for that area for Today ("Partly Cloudy" etc)
	// > Returns a characterization of whether the temperature is "hot", "cold", or "moderate"
	// (use your discretion on mapping temperatures to each type)
	weatherDetails, err = gridData.GetWeatherDetail()
	return
}

func (ws *WeatherService) getGridForecastData(api string) (data gridForeCastData, err error) {

	requestDetail := httpservice.HttpRequestDetail{
		Api: api,
	}

	var resp *httpservice.ResponseData
	resp, err = ws.client.Get(context.Background(), &requestDetail)
	if err != nil {
		return gridForeCastData{}, err
	}

	var responseBody *bytes.Reader
	responseBody, err = resp.GetBody()
	if err != nil {
		return gridForeCastData{}, err
	}

	model := new(gridForeCastData)
	err = getJson(responseBody, model)
	if err != nil {
		return gridForeCastData{}, err
	}

	return *model, nil
}
func (ws *WeatherService) getGridForecastEndPoint(api string) (endpoint string, err error) {

	requestDetail := httpservice.HttpRequestDetail{
		Api: api,
	}

	var resp *httpservice.ResponseData
	resp, err = ws.client.Get(context.Background(), &requestDetail)
	if err != nil {
		return
	}

	var responseBody *bytes.Reader
	responseBody, err = resp.GetBody()
	if err != nil {
		return
	}

	model := new(gridMetaData)
	err = getJson(responseBody, model)
	if err != nil {
		err = errorlib.NewResponseError(500, err.Error())
		return
	}

	endpoint = model.Properties.ForecastGridData
	return
}

func getJson(responseBody *bytes.Reader, target interface{}) error {
	return json.NewDecoder(responseBody).Decode(target)
}
