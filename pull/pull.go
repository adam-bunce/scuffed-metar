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
		globals.Logger.Printf("Failed to get cameco response for %s err: %v", airportCode, err)
		return camecoMetarInfo
	}

	// Read Body
	var resBody types.CamecoResponse
	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		globals.Logger.Printf("Failed to read cameco body for airport %s err: %v", airportCode, err)
		return camecoMetarInfo
	}
	defer res.Body.Close()
	err = json.Unmarshal(resBodyBytes, &resBody)
	if err != nil {
		globals.Logger.Printf("Failed to unmarshall cameco body for airport %s err: %v", airportCode, err)
		return camecoMetarInfo
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

func getHighwayData(airportMetarInfo types.MetarInfo) types.MetarInfo {

	res, err := http.Get(fmt.Sprintf("http://highways.glmobile.com/%s", airportMetarInfo.AirportName))
	if err != nil {
		globals.Logger.Printf("Failed to get highway page for airport code %s err: %v", airportMetarInfo.AirportName, err)
		return airportMetarInfo
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		globals.Logger.Printf("Failed to read highway body for airport code %s err: %v", airportMetarInfo.AirportName, err)
		airportMetarInfo.MetarInfo = []string{"Failed to read highway body for airport"}
		return airportMetarInfo
	}

	matches := highwayMetarPattern.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		metarString := strings.Trim(matches[1], "<br>")
		metarStrings := strings.Split(metarString, "<br>")
		airportMetarInfo.MetarInfo = metarStrings

	} else {
		airportMetarInfo.MetarInfo = []string{"Failed to find METAR RegEx matches"}
		return airportMetarInfo
	}

	return airportMetarInfo
}

func GetAllHighwayData() []types.MetarInfo {
	var highwayData []types.MetarInfo
	// urls are based on name but want to display airport codes
	highwayAirports := []types.MetarInfo{{
		AirportCode: "CZFD",
		AirportName: "fonddulac",
	}, {
		AirportCode: "CZWL",
		AirportName: "wollaston",
	}}
	for _, airportName := range highwayAirports {
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
		globals.Logger.Printf("Failed to get Point North CYNL Metar err: %v", err)
		return pointsNorthData
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		globals.Logger.Printf("Failed to parse Point North CYNL Metar HTML err: %v", err)
		return pointsNorthData
	}

	matches := pointsNorthRegex.FindAllStringSubmatch(string(body), -1)

	if len(matches) > 1 {
		for _, metar := range matches {
			pointsNorthData.MetarInfo = append(pointsNorthData.MetarInfo, metar[1])
		}
	}

	return pointsNorthData
}

// CYLJ, CYSF                        already exist in cams (edited)
// CYVC CYKJ CYPA         are new

func GetNavCanadaMetars() []types.MetarInfo {
	var navCanadaMetars []types.MetarInfo

	// actually god api so i can get it all in one request
	// endpoint := fmt.Sprintf("https://plan.navcanada.ca/weather/api/alpha/?point=%s&alpha=metar&alpha=taf&metar_choice=3&_=1698010681463", point)
	// &_=1698010681480 is on the end idk what it is used for though so remove dit
	// ?point=CYVT|site|-108.418,55.842&alpha=metar
	endpoint := "https://plan.navcanada.ca/weather/api/alpha/?point=CYVT|site|-108.418,55.842&point=CYLJ|site|-108.523,54.125&point=CYSF|site|-105.841,59.250&point=CYVC|site|-105.267,55.151&point=CYKJ|site|-105.617,57.256&point=CYPA|site|-105.673,53.215&alpha=metar&alpha=taf&metar_choice=3"
	res, err := http.Get(endpoint)
	if err != nil {
		globals.Logger.Printf("Failed to get nav canada metar err: %v", err)
		return navCanadaMetars
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	var navCanadaBodyData types.NavCanadaData
	err = json.Unmarshal(body, &navCanadaBodyData)
	if err != nil {
		globals.Logger.Printf("Failed to unmarshall nav canada metar err: %v", err)
		return navCanadaMetars
	}

	// most recent metar is the last for some reason
	slices.Reverse(navCanadaBodyData.Data)
	for _, data := range navCanadaBodyData.Data {
		if data.Type == "metar" || data.Type == "taf" {
			// find matching
			airportExists := false
			for i, airport := range navCanadaMetars {
				if airport.AirportCode == data.Location {
					// exists already add to data
					airportExists = true
					navCanadaMetars[i].MetarInfo = append(airport.MetarInfo, data.Text)
				}
			}
			if !airportExists {
				// update with airport
				navCanadaMetars = append(navCanadaMetars, types.MetarInfo{
					AirportCode: data.Location,
					MetarInfo:   []string{data.Text},
				})
			}

		}
	}

	return navCanadaMetars
}
