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

// ['CYZL'] -> {"CYZL", "Place Airport", "place", 2} 															// http://highways.glmobile.com/ normal cam
// ['CYZL'] -> {"CYZL", "Goober Airport", "goober", 3} 															// http://highways.glmobile.com/ normal cam
// ['CYZL'] -> {"CYZL", "Goober Airport", "goober", 3} 															// http://highways.glmobile.com/ normal cam
// ['CPZL'] -> {"CPZL", "Sandy Bay", "", 0, "http://highways.glmobile.com/",{"/ptz.jpg", "/ptz2.jpg", "/ptz3.jpg"}}		// highways outlier cam
// ['CBJL'] -> {"CBJL", "TEST Airport", "", 0, "saskpopwer.glmobile.com/charlot", {"/runway.jpg", "/hill.jpg"}}  // special cam
// ['CBJL'] -> {"CBJL", "TEST Airport"} 																		// just metar
// ['CBJL'] -> {"CBJL", "TEST Airport"} 																		// just metar
// ['CBJL'] -> {"CBJL", "TEST Airport"} 																		// just metar
// ['CBJL'] -> {"CBJL", "TEST Airport"} 																		// just metar
// ['CBJL'] -> {"CBJL", "TEST Airport"} 																		// just metar

// PULLING DATA:
// CAMECO: airportCode
// Highway: airportSpecialName (pass in airportName/Code pairs?)
// PointsNorth: 1 thing, gonna pass in airport code?
// NavCanada: pass in query parameters, or just hardcode...

// each functions returns a map of ['CYZL'] -> []string that's used as it's metar info
// the airports that each function pulls are determine by _some_ input to the function (except points north lol)
// AirportInfo being bloated is find i think, better to keep logic out of tempaltes supposedly

type AirportInfo struct {
	AirportName string // lowkey don't need
	AirportCode string

	CameraAirportIdentifier string // this is for highways.glmobile.com, pattern is <base>/airportName/ptz{int}.jpg
	CameraCount             int

	// yo i can maybe make this better and make the above two
	// not needed anymore because they all follow the same pattern?
	// kinda getting too extra thouuu
	CamHarcodedBaseUrl  string
	CameraUrlsHardCoded []string // CJP9, sandy bay

	MetarInfo []string
}

type IndexData struct {
	sync.Mutex
	AirportInformation map[string]AirportInfo
	LastUpdate         time.Time
}
