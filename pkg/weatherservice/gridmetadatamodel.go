package weatherservice

type gridMetaData struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Properties gridProperties `json:"properties"`
}

type gridProperties struct {
	Geometry            string `json:"geometry"`
	ID                  string `json:"@id"`
	Type                string `json:"@type"`
	Cwa                 string `json:"cwa"`
	ForecastOffice      string `json:"forecastOffice"`
	GridID              string `json:"gridId"`
	GridX               int    `json:"gridX"`
	GridY               int    `json:"gridY"`
	Forecast            string `json:"forecast"`
	ForecastHourly      string `json:"forecastHourly"`
	ForecastGridData    string `json:"forecastGridData"`
	ObservationStations string `json:"observationStations"`
	ForecastZone        string `json:"forecastZone"`
	County              string `json:"county"`
	FireWeatherZone     string `json:"fireWeatherZone"`
	TimeZone            string `json:"timeZone"`
	RadarStation        string `json:"radarStation"`
}
