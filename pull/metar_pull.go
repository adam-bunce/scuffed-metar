package pull

import (
	"encoding/json"
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/types"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"sync"
)

func GetAllCamecoData(dataChan chan<- types.WeatherPullInfo, wg *sync.WaitGroup) {
	for _, airportCode := range []string{"CJW7", "CYKC", "CKQ8"} {
		wg.Add(1)
		go func(ac string) { getCamecoData(ac, dataChan); wg.Done() }(airportCode)
	}
	wg.Done()
}

func getCamecoData(airportCode string, dataChan chan<- types.WeatherPullInfo) {
	weatherInfo := types.WeatherPullInfo{
		AirportCode: airportCode,
	}

	var camecoRequestBody = strings.NewReader(fmt.Sprintf(`{
	   "request": {
	       "__type": "WebDataRequest:http://COM.AXYS.COMMON.WEB.CONTRACTS",
	       "Key": "METAR",
	       "DataSourceKey": "7e7dbc35-1d26-4b85-8f7e-077ad7bad794",
	       "Query": "SELECT * FROM avWX_%s_METAR WHERE DataTimeStamp >= DATEADD(DAY, -1, GETUTCDATE()) ORDER BY DataTimeStamp DESC"
	   }
	}`, airportCode))

	req, err := http.NewRequest("POST", "http://smartweb.axys-aps.com/svc/WebDataService.svc/WebData/GetWebDataResponse", camecoRequestBody)
	if err != nil {
		globals.Logger.Printf("Failed to create cameco request for %s, err: %v", airportCode, err)
		weatherInfo.Error = err
		dataChan <- weatherInfo
		return
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Keep-Alive", "timeout=3")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	res, err := globals.Client.Do(req)
	if err != nil {
		globals.Logger.Printf("Failed to get cameco response for %s err: %v", airportCode, err)
		weatherInfo.Error = err
		dataChan <- weatherInfo
		return
	}
	defer res.Body.Close()

	var resBody types.CamecoResponse
	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		globals.Logger.Printf("Failed to read cameco body for airport %s err: %v", airportCode, err)
		weatherInfo.Error = err
		dataChan <- weatherInfo
		return
	}

	err = json.Unmarshal(resBodyBytes, &resBody)
	if err != nil {
		globals.Logger.Printf("Failed to unmarshall cameco body for airport %s err: %v", airportCode, err)
		weatherInfo.Error = err
		dataChan <- weatherInfo
		return
	}

	for i, row := range resBody.D.Rows {
		if i == 5 {
			break
		}

		metar := strings.Split(row.RowData, ",")
		if len(metar) > 1 {
			weatherInfo.Metar = append(weatherInfo.Metar, metar[1])
		} else {
			globals.Logger.Printf("Failed to parse metar: %s", airportCode)
			weatherInfo.Error = err
			dataChan <- weatherInfo
			return
		}
	}
	dataChan <- weatherInfo
}

// NOTE: this will probably change once the metar testing is done
var specialHighwayMetarPattern = regexp.MustCompile(`(?s)</h2>.*<p>(.*?)<br>(.*)<b></b>(.*)</p>`) // this is horrible
var specialHighwayTestWarning = regexp.MustCompile(`<h2>(.*?)</h2>`)                              // temp for warning

var highwayMetarPattern = regexp.MustCompile(`(?s)</h1>\s*(.*?)\s*<b>`)

func GetAllHighwayData(dataChan chan<- types.WeatherPullInfo, wg *sync.WaitGroup) {
	airportInfo := []string{"CZFD", "fonddulac", "CZWL", "wollaston"}
	for i := 0; i < len(airportInfo); i += 2 {
		wg.Add(1)
		go func(i int) {
			getHighwayData(airportInfo[i+1], airportInfo[i], dataChan)
			wg.Done()
		}(i)
	}
	//specialAirportInfo := []string{"CJY4", "sandybay", "CJL4", "laloche", "CJF3", "ilealacrosse"}
	//for i := 0; i < len(specialAirportInfo); i += 2 {
	//	wg.Add(1)
	//	go func(i int) {
	//		getSpecialHighwayData(specialAirportInfo[i+1], specialAirportInfo[i], dataChan)
	//		wg.Done()
	//	}(i)
	//}
	wg.Done()
}

func getSpecialHighwayData(airportName, airportCode string, dataChan chan<- types.WeatherPullInfo) {
	weatherInfo := types.WeatherPullInfo{
		AirportCode: airportCode,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://highways.glmobile.com/%s", airportName), nil)
	if err != nil {
		globals.Logger.Printf("Failed to create highway request for airport code %s err: %v", airportName, err)
		weatherInfo.Error = err
		dataChan <- weatherInfo
		return
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Keep-Alive", "timeout=3")

	res, err := globals.Client.Do(req)
	if err != nil {
		globals.Logger.Printf("Failed to get highway page for airport code %s err: %v", airportName, err)
		weatherInfo.Error = err
		dataChan <- weatherInfo
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		globals.Logger.Printf("Failed to read highway body for airport code %s err: %v", airportName, err)
		weatherInfo.Error = err
		dataChan <- weatherInfo
		return
	}

	warningMessageMatches := specialHighwayTestWarning.FindStringSubmatch(string(body))
	if len(warningMessageMatches) > 1 {
		weatherInfo.Error = fmt.Errorf(warningMessageMatches[1])
	}

	metarMatches := specialHighwayMetarPattern.FindStringSubmatch(string(body))
	if len(metarMatches) != 0 {
		linesInH2Para := strings.Split(metarMatches[0], "\n")
		var metarSpeciMatches []string

		for _, item := range linesInH2Para {
			if strings.Contains(item, "METAR") || strings.Contains(item, "SPECI") {
				metarSpeciMatches = append(metarSpeciMatches, item)
			}
		}

		weatherInfo.Metar = metarSpeciMatches
	} else {
		weatherInfo.Error = fmt.Errorf("failed to find metar matches")
	}

	dataChan <- weatherInfo
}

func getHighwayData(airportName, airportCode string, dataChan chan<- types.WeatherPullInfo) {
	weatherInfo := types.WeatherPullInfo{
		AirportCode: airportCode,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://highways.glmobile.com/%s", airportName), nil)
	if err != nil {
		globals.Logger.Printf("Failed to create request for highway page, airport code %s err: %v", airportName, err)
		weatherInfo.Error = err
		dataChan <- weatherInfo
		return
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Keep-Alive", "timeout=3")

	res, err := globals.Client.Do(req)
	if err != nil {
		globals.Logger.Printf("Failed to get highway page for airport code %s err: %v", airportName, err)
		weatherInfo.Error = err
		dataChan <- weatherInfo
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		globals.Logger.Printf("Failed to read highway body for airport code %s err: %v", airportName, err)
		weatherInfo.Error = err
		dataChan <- weatherInfo
		return
	}

	matches := highwayMetarPattern.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		metarString := strings.Trim(matches[1], "<br>")
		metarStrings := strings.Split(metarString, "<br>")
		weatherInfo.Metar = metarStrings
	} else {
		weatherInfo.Error = fmt.Errorf("failed to find metar matches")
	}

	dataChan <- weatherInfo
}

var pointsNorthRegex = regexp.MustCompile(`(?i)<TD COLSPAN="3">(.*?)</TD>`)

func GetPointsNorthMetar(dataChan chan<- types.WeatherPullInfo, wg *sync.WaitGroup) {
	pointsNorthData := types.WeatherPullInfo{AirportCode: "CYNL"}

	// only one we care about so hard coding it for now
	res, err := globals.Client.Get("https://www.pointsnorthgroup.ca/weather/CYNL_metar.html")
	if err != nil {
		globals.Logger.Printf("Failed to get Point North CYNL Metar err: %v", err)
		pointsNorthData.Error = err
		dataChan <- pointsNorthData
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		globals.Logger.Printf("Failed to parse Point North CYNL Metar HTML err: %v", err)
		pointsNorthData.Error = err
		dataChan <- pointsNorthData
		return
	}

	matches := pointsNorthRegex.FindAllStringSubmatch(string(body), -1)

	if len(matches) > 1 {
		for _, metar := range matches {
			pointsNorthData.Metar = append(pointsNorthData.Metar, metar[1])
		}
	}

	dataChan <- pointsNorthData
	wg.Done()
}

func GetNavCanadaMetars(dataChan chan<- types.WeatherPullInfo, wg *sync.WaitGroup) {
	navCanadaMetars := make(map[string]types.WeatherPullInfo)

	endpoint := "https://plan.navcanada.ca/weather/api/alpha/?" +
		"site=CYXE&" +
		"site=CYVT&" +
		"site=CYLJ&" +
		"site=CYSF&" +
		"site=CYVC&" +
		"site=CYKJ&" +
		"site=CYPA&" +
		"alpha=metar&" +
		"alpha=taf&" +
		"metar_choice=3"

	res, err := globals.Client.Get(endpoint)
	if err != nil {
		globals.Logger.Printf("Failed to get nav canada metar err: %v", err)
		for k := range navCanadaMetars {
			airport := navCanadaMetars[k]
			airport.Error = err
			dataChan <- navCanadaMetars[k]
		}
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	var navCanadaResp types.NavCanadaResponse
	err = json.Unmarshal(body, &navCanadaResp)
	if err != nil {
		globals.Logger.Printf("Failed to unmarshall nav canada metar err: %v", err)
		for k := range navCanadaMetars {
			airport := navCanadaMetars[k]
			airport.Error = err
			dataChan <- navCanadaMetars[k]
		}
		return
	}

	// most recent metar is the last in order
	slices.Reverse(navCanadaResp.Data)
	for _, data := range navCanadaResp.Data {
		airport, ok := navCanadaMetars[data.Location]
		if !ok {
			navCanadaMetars[data.Location] = types.WeatherPullInfo{AirportCode: data.Location}
			airport = navCanadaMetars[data.Location]
		}

		airport.AirportCode = data.Location
		if data.Type == "metar" {
			airport.Metar = append(airport.Metar, data.Text)
		} else if data.Type == "taf" {
			airport.Taf = append(airport.Taf, data.Text)
		}
		navCanadaMetars[data.Location] = airport

	}

	for _, v := range navCanadaMetars {
		dataChan <- v
	}
	wg.Done()
}
