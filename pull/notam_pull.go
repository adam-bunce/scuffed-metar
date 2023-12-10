package pull

import (
	"encoding/json"
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/types"
	"sort"
	"strings"
	"time"
)

func GetNotam(airportCodes []string) ([]types.NotamData, error) {
	endpoint := "https://plan.navcanada.ca/weather/api/alpha/?alpha=notam&notam_choice=english"
	for _, airport := range airportCodes {
		endpoint += "&site=" + airport
	}

	notams := make(map[string]types.NotamData)
	var notamSlice []types.NotamData

	response, err := globals.Client.Get(endpoint)
	if err != nil {
		globals.Logger.Printf("failed to get NOTAMS", err)
		return notamSlice, fmt.Errorf("failed to get NOTAMS")
	}
	var bodyValue types.NavCanadaResponse
	err = json.NewDecoder(response.Body).Decode(&bodyValue)
	if err != nil {
		globals.Logger.Printf("failed to decode response body", err)
		return notamSlice, fmt.Errorf("failed to decode response body")
	}

	for _, item := range bodyValue.Data {
		// .Text is escaped json
		var parsedText types.NotamParsedText
		err = json.NewDecoder(strings.NewReader(item.Text)).Decode(&parsedText)
		if err != nil {
			globals.Logger.Printf("failed to decode response text field", err)
			return notamSlice, fmt.Errorf("failed to decode response text field")
		}
		applicableAirportCodes := append(notams[parsedText.Raw].ApplicableAirports, item.Position.PointReference)

		var startValidityTime time.Time
		var endValidityTime time.Time
		if item.StartValidity != "" {
			startValidityTime, err = time.Parse(types.NavCanadaTimeFormat, item.StartValidity)
			if err != nil {
				globals.Logger.Printf("failed to parse start validity", err)
				return notamSlice, fmt.Errorf("failed to parse start validity")
			}
		}
		if item.EndValidity != "" {
			endValidityTime, err = time.Parse(types.NavCanadaTimeFormat, item.EndValidity)
			if err != nil {
				globals.Logger.Printf("failed to parse end validity", err)
				return notamSlice, fmt.Errorf("failed to parse end validity")
			}
		}

		notams[parsedText.Raw] = types.NotamData{
			ApplicableAirports: applicableAirportCodes,
			StartValidity:      startValidityTime,
			EndValidity:        endValidityTime,
			Notam:              parsedText.Raw,
		}

	}

	// sort by start validity descending
	for _, v := range notams {
		notamSlice = append(notamSlice, v)
	}
	sort.Slice(notamSlice, func(a, b int) bool {
		return notamSlice[a].StartValidity.After(notamSlice[b].StartValidity)
	})

	return notamSlice, nil
}
