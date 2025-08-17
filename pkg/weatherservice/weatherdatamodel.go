package weatherservice

type gridForeCastData struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Properties gridForecastProperties `json:"properties"`
}

type maxTempValues struct {
	ValidTime string  `json:"validTime"`
	Value     float64 `json:"value"`
}

type maxTemperatureData struct {
	Uom    string          `json:"uom"`
	Values []maxTempValues `json:"values"`
}

type gridForecastProperties struct {
	ID             string             `json:"@id"`
	Type           string             `json:"@type"`
	UpdateTime     string             `json:"updateTime"`
	ValidTimes     string             `json:"validTimes"`
	ForecastOffice string             `json:"forecastOffice"`
	GridID         string             `json:"gridId"`
	GridX          int                `json:"gridX"`
	GridY          int                `json:"gridY"`
	MaxTemperature maxTemperatureData `json:"maxTemperature"`
	Weather        WeatherData        `json:"weather"`
}

type WeatherData struct {
	Values []WeatherValueData `json:"values"`
}

type WeatherValueData struct {
	ValidTime string         `json:"validTime"`
	Value     []WeatherValue `json:"value"`
}

type WeatherValue struct {
	Coverage   string         `json:"coverage"`
	Weather    *string        `json:"weather"`
	Intensity  string         `json:"intensity"`
	Visibility VisibilityData `json:"visibility"`
	Attributes []string       `json:"attributes"`
}

type VisibilityData struct {
	Value          int    `json:"value"`
	MaxValue       int    `json:"maxValue"`
	MinValue       int    `json:"minValue"`
	UnitCode       string `json:"unitCode"`
	QualityControl string `json:"qualityControl"`
}

func (data gridForeCastData) GetWeatherDetail() (WeatherDetails, error) {
	currentDate := getCurrentDate()

	var temperature float64
	var weather string
	for _, mexTempvalue := range data.Properties.MaxTemperature.Values {
		parsedDate := mexTempvalue.ValidTime[0:10]
		if parsedDate == currentDate {
			temperature = mexTempvalue.Value
			break
		}
	}
	for _, weatherValueData := range data.Properties.Weather.Values {
		parsedDate := weatherValueData.ValidTime[0:10]
		if parsedDate == currentDate {
			weather = findWeather(weatherValueData.Value)
			weather = "weather on " + parsedDate + " is:" + weather
			break
		}
	}
	return WeatherDetails{
		Temperature: temperature,
		Weather:     weather,
	}, nil
}

func findWeather(weatherData []WeatherValue) string {
	for _, value := range weatherData {
		if value.Weather != nil {
			return *value.Weather
		}
	}
	return "unknown, because data.Properties.Weather.Values[currentDate] received from api.weather.gov is null"
}
