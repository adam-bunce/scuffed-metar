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

type NavCanadaImageResponse struct {
	Meta struct {
		Now   string `json:"now"`
		Count struct {
			Image int `json:"image"`
		} `json:"count"`
		Messages []interface{} `json:"messages"`
	} `json:"meta"`
	Data []struct {
		Type          string      `json:"type"`
		Pk            string      `json:"pk"`
		Location      string      `json:"location"`
		StartValidity string      `json:"startValidity"`
		EndValidity   interface{} `json:"endValidity"`
		Text          string      `json:"text"`
		HasError      bool        `json:"hasError"`
		Position      struct {
			PointReference string `json:"pointReference"`
			RadialDistance int    `json:"radialDistance"`
		} `json:"position"`
	} `json:"data"`
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
}

type WeatherPullInfo struct {
	AirportCode string
	WeatherInfo
}

type IndexData struct {
	sync.Mutex

	Version            string
	AirportInformation []AirportInfo
	LastUpdate         time.Time
}

type GfaPageData struct {
	sync.Mutex

	Version string
	GfaInfo
	LastUpdate time.Time

	Error error
}

type NotamPageData struct {
	Version      string
	NoTamOptions []string

	Raw        string `json:"raw"`
	English    string `json:"english"`
	French     string `json:"french"`
	LastUpdate time.Time
}

type NotamParsedText struct {
	Raw     string `json:"raw"`
	English string `json:"english"`
	French  string `json:"french"`
}
