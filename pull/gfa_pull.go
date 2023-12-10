package pull

import (
	"encoding/json"
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/types"
	"strconv"
	"strings"
)

// GetGFAImageIds gets the most recently issued forecast images for ice/turbulence/freezing level and clouds/weather
// data return from endpoint has image id's in .Text, .Text consists of 3 FrameLists the third FrameList
// has the most recently issued images, each .Frame has one image
//
// Data _should_ have a length of 2 the CloudsWeather or IcingTurbFreezing are distinguished by the .Location
// which can contain CLDWX or TURBC
func GetGFAImageIds() (types.GfaInfo, error) {
	var imageIds types.GfaInfo

	response, err := globals.Client.Get("https://plan.navcanada.ca/weather/api/alpha/?site=CYXE&image=GFA/CLDWX&image=GFA/TURBC")
	if err != nil {
		return imageIds, fmt.Errorf("failed to get nav canada GFA image data")
	}
	var bodyValue types.NavCanadaResponse
	err = json.NewDecoder(response.Body).Decode(&bodyValue)
	if err != nil {
		return imageIds, fmt.Errorf("failed to decode response body")
	}

	for _, item := range bodyValue.Data {
		// .Text is escaped json
		var parsedText types.ParsedText
		err = json.NewDecoder(strings.NewReader(item.Text)).Decode(&parsedText)
		if err != nil {
			return imageIds, fmt.Errorf("failed to decode response text field")
		}
		if len(parsedText.FrameLists) == 0 {
			return imageIds, fmt.Errorf("didn't receive any FrameLists")
		}
		// extract image id's from last frameset
		for _, frame := range parsedText.FrameLists[len(parsedText.FrameLists)-1].Frames {
			if len(frame.Images) > 0 {
				if strings.Contains(item.Location, "CLDWX") {
					imageIds.CloudsWeather = append(imageIds.CloudsWeather, types.GFAMetadata{
						StartValidity: frame.StartValidity,
						EndValidity:   frame.EndValidity,
						Id:            strconv.Itoa(frame.Images[len(frame.Images)-1].Id),
					})
				} else {
					imageIds.IcingTurbFreezing = append(imageIds.IcingTurbFreezing, types.GFAMetadata{
						StartValidity: frame.StartValidity,
						EndValidity:   frame.EndValidity,
						Id:            strconv.Itoa(frame.Images[len(frame.Images)-1].Id),
					})
				}
			}

		}
	}

	return imageIds, nil
}
