package pull

import (
	"encoding/json"
	"fmt"
	"github.com/adam-bunce/scuffed-metar/types"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func getCamecoData(airportCode string) types.MetarInfo {
	camecoMetarInfo := types.MetarInfo{AirportCode: airportCode}

	// Prep/Do Request
	var camecoRequestBody = strings.NewReader(fmt.Sprintf(`{
    "request": {
        "__type": "WebDataRequest:http://COM.AXYS.COMMON.WEB.CONTRACTS",
        "Key": "METAR",
        "DataSourceKey": "7e7dbc35-1d26-4b85-8f7e-077ad7bad794",
        "Query": "SELECT * FROM avWX_%s_METAR WHERE DataTimeStamp >= DATEADD(DAY, -1, GETUTCDATE()) ORDER BY DataTimeStamp DESC"
    }
}`, airportCode))

	res, err := http.Post("http://smartweb.axys-aps.com/svc/WebDataService.svc/WebData/GetWebDataResponse",
		"application/json; charset=UTF-8",
		camecoRequestBody)
	if err != nil {
		log.Printf("Failed to get cameco response for %s err: %v", airportCode, err)
		return camecoMetarInfo
	}

	// Read Body
	var resBody types.CamecoResponse
	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read cameco body for airport %s err: %v", airportCode, err)
	}
	defer res.Body.Close()
	err = json.Unmarshal(resBodyBytes, &resBody)
	if err != nil {
		log.Printf("Failed to unmarshall cameco body for airport %s err: %v", airportCode, err)
	}

	// Easy json processing, take first 5 rows
	count := 0
	for _, row := range resBody.D.Rows {
		count++
		if count == 5 {
			break
		}

		metarInfo := strings.Split(row.RowData, ",")
		if len(metarInfo) > 1 {
			camecoMetarInfo.MetarInfo = append(camecoMetarInfo.MetarInfo, metarInfo[1])
		} else {
			camecoMetarInfo.MetarInfo = append(camecoMetarInfo.MetarInfo, "failed to parse METAR info")
		}
	}

	return camecoMetarInfo
}

func GetAllCamecoData() []types.MetarInfo {
	var camecoData []types.MetarInfo
	camecoAirpotCodes := []string{"CJW7", "CYKC", "CKQ8"}
	for _, airportCode := range camecoAirpotCodes {
		camecoData = append(camecoData, getCamecoData(airportCode))
	}

	return camecoData
}

var highwayMetarPattern = regexp.MustCompile(`(?s)</h1>\s*(.*?)\s*<b>`)

func getHighwayData(airportName string) types.MetarInfo {
	highwayData := types.MetarInfo{
		AirportName: airportName,
	}

	res, err := http.Get(fmt.Sprintf("http://highways.glmobile.com/%s", airportName))
	if err != nil {
		log.Printf("Failed to get highway page for airport code %s err: %v", airportName, err)
		return highwayData
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read highway body for airport code %s err: %v", airportName, err)
		highwayData.MetarInfo = []string{"Failed to read highway body for airport"}
		return highwayData
	}

	matches := highwayMetarPattern.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		metarString := strings.Trim(matches[1], "<br>")
		metarStrings := strings.Split(metarString, "<br>")
		highwayData.MetarInfo = metarStrings

	} else {
		highwayData.MetarInfo = []string{"Failed to find METAR RegEx matches"}
		return highwayData
	}

	return highwayData
}

func GetAllHighwayData() []types.MetarInfo {
	var highwayData []types.MetarInfo
	highwayAirportNames := []string{"Fonddulac", "Wollaston"}
	for _, airportName := range highwayAirportNames {
		highwayData = append(highwayData, getHighwayData(airportName))
	}

	return highwayData
}

var pointsNorthRegex = regexp.MustCompile(`(?i)<TD COLSPAN="3">(.*?)</TD>`)

func GetPointsNorthMetar() types.MetarInfo {
	pointsNorthData := types.MetarInfo{AirportCode: "CYNL"}
	// only one we care about so hard coding it for now
	res, err := http.Get("https://www.pointsnorthgroup.ca/weather/CYNL_metar.html")
	if err != nil {
		log.Printf("Failed to get Point North CYNL Metar err: %v", err)
		return pointsNorthData
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to parse Point North CYNL Metar HTML err: %v", err)
	}

	matches := pointsNorthRegex.FindAllStringSubmatch(string(body), -1)

	if len(matches) > 1 {
		for _, metar := range matches {
			pointsNorthData.MetarInfo = append(pointsNorthData.MetarInfo, metar[1])
		}
	}

	return pointsNorthData
}
