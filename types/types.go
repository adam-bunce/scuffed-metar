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
	MetarData  []MetarInfo
	Cameras    []WxCam
	LastUpdate time.Time
}
