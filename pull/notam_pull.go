package pull

import (
	"encoding/json"
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/types"
	"sort"
	"strings"
	"sync"
	"time"
)

func GetAllNotams(airportCodes []string) []types.NotamData {
	// NavCan API supports getting notams for multiple sites BUT, if i add radius=n the response format changes
	// making it <?impossible/really hard?> to figure out what site is associated with which notam's so i do every one individually
	var wg sync.WaitGroup
	respChan := make(chan types.NavCanadaResponse)
	for _, airport := range airportCodes {
		wg.Add(1)
		go func(airport string) { GetSingleNotam(airport, respChan); wg.Done() }(airport)
	}

	go func() {
		wg.Wait()
		close(respChan)
	}()

	notams := make(map[string]types.NotamData)
	for data := range respChan {
		for _, notam := range data.Data {
			var parsedText types.NotamParsedText
			err := json.NewDecoder(strings.NewReader(notam.Text)).Decode(&parsedText)
			if err != nil {
				globals.Logger.Printf("failed to decode response text field", err)
			}

			var startValidityTime time.Time
			var endValidityTime time.Time
			if notam.StartValidity != "" {
				startValidityTime, err = time.Parse(types.NavCanadaTimeFormat, notam.StartValidity)
				if err != nil {
					globals.Logger.Printf("failed to parse start validity", err)
				}
			}
			if notam.EndValidity != "" {
				endValidityTime, err = time.Parse(types.NavCanadaTimeFormat, notam.EndValidity)
				if err != nil {
					globals.Logger.Printf("failed to parse end validity", err)
				}
			}

			applicableAirportCodes := append(notams[parsedText.Raw].ApplicableAirports, notam.Position.PointReference)

			notams[parsedText.Raw] = types.NotamData{
				ApplicableAirports: applicableAirportCodes,
				StartValidity:      startValidityTime,
				EndValidity:        endValidityTime,
				Notam:              parsedText.Raw,
			}
		}
	}

	var notamSlice []types.NotamData

	// sort by start validity descending
	for _, v := range notams {
		notamSlice = append(notamSlice, v)
	}
	sort.Slice(notamSlice, func(a, b int) bool {
		return notamSlice[a].StartValidity.After(notamSlice[b].StartValidity)
	})

	return notamSlice
}

func GetSingleNotam(airportCode string, dataChan chan<- types.NavCanadaResponse) {
	endpoint := "https://plan.navcanada.ca/weather/api/alpha/?alpha=notam&radius=10&notam_choice=english&site=" + airportCode

	response, err := globals.Client.Get(endpoint)
	if err != nil {
		globals.Logger.Printf("failed to get NOTAMS %v", err)
		globals.SendWebhook(fmt.Sprintf("[NOTAM] get request failed err: %v, endpoint: %s ", err, endpoint))
		dataChan <- types.NavCanadaResponse{}
		return
	}

	var bodyValue types.NavCanadaResponse
	err = json.NewDecoder(response.Body).Decode(&bodyValue)
	if err != nil {
		globals.Logger.Printf("failed to decode response body %v", err)
		globals.SendWebhook(fmt.Sprintf("[NOTAM] [types.NavCanadaResponse] failed to decode response body: err: %v, endpoint: %s ", err, endpoint))
		dataChan <- types.NavCanadaResponse{}
		return
	}

	dataChan <- bodyValue
}
