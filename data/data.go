package data

type TimeSeriesData struct {
	DateTime  string `json:"datetime"`
	Value     string `json:"value"`
	Partition string `json:"partition"`
}
