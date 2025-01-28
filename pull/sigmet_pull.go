package pull

import (
	"encoding/json"
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/types"
	"sync"
)

type SigmetInfo struct {
	Site string
	Text string
}

func GetAllSigmetAirmet(airports []string) map[string][]string {
	respChan := make(chan SigmetInfo)
	var wg sync.WaitGroup

	for _, airport := range airports {
		wg.Add(1)
		go func(airport string) { GetSigmetAirmet(airport, respChan); wg.Done() }(airport)
	}

	go func() {
		wg.Wait()
		close(respChan)
	}()

	resSet := make(map[string][]string) // sigmet -> airportCode
	for text := range respChan {
		resSet[text.Text] = append(resSet[text.Text], text.Site)
	}

	return resSet
}

func GetSigmetAirmet(airportCode string, respChan chan SigmetInfo) error {
	radius := 50
	endpoint := fmt.Sprintf("https://plan.navcanada.ca/weather/api/alpha/?alpha=sigmet&alpha=airmet&radius=%d&site=%s", radius, airportCode)

	response, err := globals.Client.Get(endpoint)
	if err != nil {
		globals.Logger.Printf("%v", err)
		return err
	}
	defer response.Body.Close()

	var bodyValue types.NavCanadaResponse
	err = json.NewDecoder(response.Body).Decode(&bodyValue)
	if err != nil {
		globals.Logger.Printf("%v", err)
	}

	if len(bodyValue.Data) == 0 {
		respChan <- SigmetInfo{Site: airportCode, Text: "None"}
		return nil
	}

	for _, text := range bodyValue.Data {
		respChan <- SigmetInfo{Site: airportCode, Text: text.Text}
	}

	return nil
}
