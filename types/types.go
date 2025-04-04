package types

import (
	"encoding/xml"
	"sync"
	"time"
)

const NavCanadaTimeFormat = "2006-01-02T15:04:05"
const NavCanadaTimeFormatAlt = "2006-01-02T15:04:05+00:00"

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
		Type          string `json:"type"`
		Pk            string `json:"pk"`
		Location      string `json:"location"`
		StartValidity string `json:"startValidity"`
		EndValidity   string `json:"endValidity"`
		Text          string `json:"text"`
		HasError      bool   `json:"hasError"`
		Position      struct {
			PointReference string `json:"pointReference"`
			RadialDistance int    `json:"radialDistance"`
		} `json:"position"`
	} `json:"data"`
}

type MesotechResponse struct {
	XMLName   xml.Name    `xml:"AWA"`
	Copyright string      `xml:"copyright,attr"`
	SiteName  string      `xml:"SITE_NAME"`
	Current   interface{} `xml:"CURRENT"`
	ReportLog []struct {
		ReportLogId   string `xml:"ReportLogId"`
		DateTimeStamp string `xml:"DateTimeStamp"`
		ReportType    string `xml:"ReportType"`
		Report        string `xml:"Report"`
	} `xml:"ReportLog>Row"`
}

type GfaInfo struct {
	CloudsWeather     []GFAMetadata
	IcingTurbFreezing []GFAMetadata
}

type GFAMetadata struct {
	StartValidity string
	EndValidity   string
	Id            string
}

type ParsedText struct {
	Product      string `json:"product"`
	SubProduct   string `json:"sub_product"`
	Geography    string `json:"geography"`
	SubGeography string `json:"sub_geography"`
	FrameLists   []struct {
		Id     int    `json:"id"`
		Sv     string `json:"sv"`
		Ev     string `json:"ev"`
		Frames []struct {
			Id            int    `json:"id"`
			StartValidity string `json:"sv"`
			EndValidity   string `json:"ev"`
			Images        []struct {
				Id      int    `json:"id"`
				Created string `json:"created"`
			} `json:"images"`
		} `json:"frames"`
	} `json:"frame_lists"`
}

type WeatherInfo struct {
	Metar []string
	Taf   []string
	Error error
}

type AirportInfo struct {
	AirportName string
	AirportCode string

	CamBaseUrl string
	CamPicUrls []string
	WeatherInfo

	MetarSource   string
	AwosFrequency string

	Stale int
}

type WeatherPullInfo struct {
	AirportCode      string
	UpdatedImageUrls *[]string
	WeatherInfo
}

type GenericPageData struct {
	LastUpdate time.Time
	Version    string
}

type IndexData struct {
	sync.Mutex

	AirportInformation []AirportInfo
	GenericPageData
	CamecoLastUpdated time.Time
}

type GfaPageData struct {
	sync.Mutex

	GenericPageData
	GfaInfo

	Error error
}

type NotamPageData struct {
	GenericPageData
	NoTamOptions []string
	Error        error
}

type NotamParsedText struct {
	Raw     string `json:"raw"`
	English string `json:"english"`
	French  string `json:"french"`
}

type NotamData struct {
	ApplicableAirports []string
	StartValidity      time.Time
	EndValidity        time.Time
	Notam              string
}

type WindsPageData struct {
	GenericPageData
	WindsOptions []string
	Error        error
	MaxInt       float64
}

type WindsData struct {
	AirportCode string

	High []WindData
	Low  []WindData

	MaxInt float64
}

type WindData struct {
	Data []ElevationValueCombo

	BasedOn     time.Time
	Valid       time.Time
	ForUseStart time.Time
	ForUseEnd   time.Time
}

type ElevationValueCombo struct {
	Elevation float64
	Values    []float64
}

type MQTTReportLogTopicMessage struct {
	History []string `json:"history"`
}
