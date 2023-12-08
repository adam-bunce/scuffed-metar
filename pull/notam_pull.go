package pull

import (
	"encoding/json"
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/types"
	"strings"
)

func GetNotam(airportCodes []string) (map[string][]string, error) {
	endpoint := "https://plan.navcanada.ca/weather/api/alpha/?alpha=notam&notam_choice=english"
	for _, airport := range airportCodes {
		endpoint += "&site=" + airport
	}

	// string to airports
	stringsWeGot := make(map[string][]string)

	response, err := globals.Client.Get(endpoint)
	if err != nil {
		return stringsWeGot, fmt.Errorf("failed to get nav canada GFA image data")
	}
	var bodyValue types.NavCanadaImageResponse
	err = json.NewDecoder(response.Body).Decode(&bodyValue)
	if err != nil {
		return stringsWeGot, fmt.Errorf("failed to decode response body")
	}

	for _, item := range bodyValue.Data {
		// .Text is escaped json
		var parsedText types.NotamParsedText
		err = json.NewDecoder(strings.NewReader(item.Text)).Decode(&parsedText)
		if err != nil {
			return stringsWeGot, fmt.Errorf("failed to decode response text field")
		}

		stringsWeGot[parsedText.Raw] = append(stringsWeGot[parsedText.Raw], item.Position.PointReference)
	}

	return stringsWeGot, nil
}
