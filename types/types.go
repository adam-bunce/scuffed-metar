package types

import (
	"sync"
	"time"
)

type CamecoResponse struct {
	D struct {
		Type             string        `json:"__type"`
		AccessType       interface{}   `json:"AccessType"`
		Key              string        `json:"Key"`
		ModifyDateString string        `json:"ModifyDateString"`
		ModifyUser       interface{}   `json:"ModifyUser"`
		Properties       []interface{} `json:"Properties"`
		ColumnCount      int           `json:"ColumnCount"`
		Columns          []struct {
			Type       string `json:"__type"`
			ColumnName string `json:"ColumnName"`
			DataType   string `json:"DataType"`
			Ordinal    int    `json:"Ordinal"`
		} `json:"Columns"`
		RowCount int `json:"RowCount"`
		Rows     []struct {
			Type    string `json:"__type"`
			RowData string `json:"RowData"`
			RowID   int    `json:"RowID"`
		} `json:"Rows"`
	} `json:"d"`
}

type NavCanadaResponse struct {
	Meta struct {
		Now   string `json:"now"`
		Count struct {
			Metar int `json:"metar"`
			Taf   int `json:"taf"`
		} `json:"count"`
		Messages []interface{} `json:"messages"`
	} `json:"meta"`
	Data []struct {
		Type          string  `json:"type"`
		Pk            string  `json:"pk"`
		Location      string  `json:"location"`
		StartValidity string  `json:"startValidity"`
		EndValidity   *string `json:"endValidity"`
		Text          string  `json:"text"`
		HasError      bool    `json:"hasError"`
		Position      struct {
			PointReference string `json:"pointReference"`
			RadialDistance int    `json:"radialDistance"`
		} `json:"position"`
	} `json:"data"`
}

type WeatherInfo struct {
	Metar []string
	Taf   []string
}

type AirportInfo struct {
	AirportName string
	AirportCode string

	CamBaseUrl string
	CamPicUrls []string
	WeatherInfo
}

type WeatherPullInfo struct {
	AirportCode string
	WeatherInfo
	Error error
}

type IndexData struct {
	sync.Mutex

	AirportInformation []AirportInfo
	LastUpdate         time.Time
}
