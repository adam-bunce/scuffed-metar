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

type NavCanadaData struct {
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

type MetarInfo struct {
	AirportCode string
	AirportName string
	MetarInfo   []string
}

// WxCam only highway airports have cameras so all the links are formatted the same
type WxCam struct {
	AirportName string
	AirportCode string
	ImageCount  int
}

type IndexData struct {
	sync.Mutex
	MetarData  []MetarInfo // TODO need way to add urls to this so generating tempalte is cleaner
	Cameras    []WxCam     // cameras are a separate thing rn but shouldnt be
	LastUpdate time.Time
}
