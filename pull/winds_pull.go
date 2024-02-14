package pull

import (
	"encoding/json"
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/types"
	"math"
	"sort"
	"strings"
	"time"
)

// number
const (
	dataArrayIndex         = 11
	windsArrElevationIndex = 0
	maxLowWindValue        = 18_000
)

// strings
const (
	firstDateIndex   = 2
	basedOnIndex     = 3
	validIndex       = 4
	forUseStartIndex = 5
	forUseEndIndex   = 6
)

func GetWinds(airportCodes []string) ([]types.WindsData, error) {
	endpoint := "https://plan.navcanada.ca/weather/api/alpha/?alpha=upperwind&upperwind_choice=both"
	for _, airport := range airportCodes {
		endpoint += "&site=" + airport
	}

	var res []types.WindsData
	resSet := make(map[string]types.WindsData) // airportCode -> WindData (high/low)

	response, err := globals.Client.Get(endpoint)
	if err != nil {
		globals.Logger.Printf("failed to get winds %v", err)
		return res, fmt.Errorf("failed to get winds for %v", airportCodes)
	}
	defer response.Body.Close()
	var bodyValue types.NavCanadaResponse
	err = json.NewDecoder(response.Body).Decode(&bodyValue)
	if err != nil {
		globals.Logger.Printf("failed to decode response body %v", err)
		return res, fmt.Errorf("failed to decode response body")
	}

	for _, item := range bodyValue.Data {
		var (
			stringValues []string
			numberValues [][]float64
			parsedText   []interface{}
		)
		err = json.NewDecoder(strings.NewReader(item.Text)).Decode(&parsedText)
		if err != nil {
			globals.Logger.Printf("failed to decode response text field,  %v", err)
			return res, fmt.Errorf("failed to decode response text field")
		}

		// yeah...
		for i, stringVal := range parsedText {
			if i == dataArrayIndex {
				if innerArray, ok := stringVal.([]interface{}); ok {
					for _, innerItem := range innerArray {
						var tempArr []float64
						for _, arrayItem := range innerItem.([]interface{}) {
							if arrayItem != nil {
								tempArr = append(tempArr, arrayItem.(float64))
							} else {
								tempArr = append(tempArr, math.MaxInt)

							}
						}
						numberValues = append(numberValues, tempArr)
					}
				}
			} else if stringVal != nil {
				stringValues = append(stringValues, stringVal.(string))
			}
		}

		currentAirportCode := item.Position.PointReference
		winds, isHighWinds := ProcessData(stringValues, numberValues)
		if isHighWinds {
			resSet[currentAirportCode] = types.WindsData{
				AirportCode: item.Position.PointReference,
				High:        append(resSet[currentAirportCode].High, winds),
				Low:         resSet[currentAirportCode].Low,
			}
		} else {
			resSet[currentAirportCode] = types.WindsData{
				AirportCode: item.Position.PointReference,
				High:        resSet[currentAirportCode].High,
				Low:         append(resSet[currentAirportCode].Low, winds),
			}
		}
		// sort for ui order
		sort.Slice(resSet[currentAirportCode].Low, func(i, j int) bool {
			return resSet[currentAirportCode].Low[i].ForUseStart.Before(resSet[currentAirportCode].Low[j].ForUseStart)
		})
		sort.Slice(resSet[currentAirportCode].High, func(i, j int) bool {
			return resSet[currentAirportCode].High[i].ForUseStart.Before(resSet[currentAirportCode].High[j].ForUseStart)
		})
	}

	for _, v := range resSet {
		res = append(res, v)
	}

	sort.Slice(res, func(i, j int) bool {
		return strings.Compare(res[i].AirportCode, res[j].AirportCode) < 0
	})

	return res, nil
}

func ProcessData(stringValues []string, numberValues [][]float64) (types.WindData, bool) {
	var winds types.WindData
	var isHighWinds bool

	// sort based on height value
	sort.Slice(numberValues, func(i, j int) bool {
		return numberValues[i][windsArrElevationIndex] < numberValues[j][windsArrElevationIndex]
	})

	var data []types.ElevationValueCombo
	for _, windsArr := range numberValues {
		if windsArr[windsArrElevationIndex] > maxLowWindValue {
			isHighWinds = true
		}
		data = append(data, types.ElevationValueCombo{
			Elevation: windsArr[windsArrElevationIndex],
			Values:    windsArr[windsArrElevationIndex+1 : len(windsArr)-1], // last value is always zero so we remove that
		})
	}

	winds.Data = data
	winds.BasedOn, _ = time.Parse(types.NavCanadaTimeFormatAlt, stringValues[basedOnIndex])
	winds.Valid, _ = time.Parse(types.NavCanadaTimeFormatAlt, stringValues[validIndex])
	winds.ForUseStart, _ = time.Parse(types.NavCanadaTimeFormatAlt, stringValues[forUseStartIndex])
	winds.ForUseEnd, _ = time.Parse(types.NavCanadaTimeFormatAlt, stringValues[forUseEndIndex])

	return winds, isHighWinds
}
