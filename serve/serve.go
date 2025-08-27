package serve

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/pull"
	"github.com/adam-bunce/scuffed-metar/types"
	"html/template"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TODO: semantic compression for all these

//go:embed static
var files embed.FS
var fileServer = http.FileServer(http.FS(files))

//go:embed pages/waas.html
var waasTemplateString string
var waasTemplate = LoadTemplate("", "waas", waasTemplateString)

//go:embed pages/mets.html
var metsTemplateString string
var metsTemplate = LoadTemplate("", "mets", metsTemplateString)

//go:embed pages/index.html
var indexTemplateString string
var indexTemplate = LoadTemplate("", "index", indexTemplateString)

//go:embed pages/gfa.html
var gfaTemplateString string
var gfaTemplate = LoadTemplate("", "gfa", gfaTemplateString)

//go:embed pages/notam.html
var notamTemplateString string
var notamTemplate = LoadTemplate("", "notam", notamTemplateString)

//go:embed pages/winds.html
var windsTemplateString string
var windsTemplate = LoadTemplate("", "winds", windsTemplateString)

//go:embed pages/info.html
var infoTemplateString string
var infoTemplate = LoadTemplate("", "info", infoTemplateString)

//go:embed pages/trip.html
var tripTemplateString string
var tripTemplate = LoadTemplate("", "trip", tripTemplateString)

func hc(num int) []string {
	var res []string
	for i := 1; i < num+1; i++ {
		res = append(res, fmt.Sprintf("/ptz%s.jpg", strconv.Itoa(i)))
	}
	return res
}

const (
	NavCanda    = "https://plan.navcanada.ca/wxrecall/"
	Cameco      = "http://smartweb.axys-aps.com/sites/1517/"
	PointsNorth = "https://www.pointsnorthgroup.ca/weather/"
	Highways    = "http://highways.glmobile.com/"
	MetCam      = "https://www.metcam.navcanada.ca/dawc_images/wxcam/"
)

// this is the order displayed in the UI
var airportInfo = []types.AirportInfo{
	{"Cigar Lake", "CJW7", "", nil, types.WeatherInfo{}, Cameco + "CJW7/", "", 0},
	{"McArthur River", "CKQ8", "", nil, types.WeatherInfo{}, Cameco + "CKQ8/", "", 0},
	{"Collins Bay / Rabbit Lake", "CYKC", "", nil, types.WeatherInfo{}, Cameco + "CYKC/", "", 0},
	{"Key Lake", "CYKJ", MetCam + "CYKJ", []string{"/CYKJ_SW-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, "", 0},

	{"Points North", "CYNL", "https://pointsnorthgroup.ca/liveview", []string{"/Northview.jpg", "/Southview.jpg"}, types.WeatherInfo{}, PointsNorth + "CYNL_metar.html", "", 0},

	{"Saskatoon (Stoon)", "CYXE", "", nil, types.WeatherInfo{}, NavCanda, "", 0},
	{"Prince Albert (PA)", "CYPA", "https://www.metcam.navcanada.ca/dawc_images/wxcam/CYPA", []string{"/CYPA_E-full-e.jpeg", "/CYPA_N-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, "", 0},
	{"La Ronge", "CYVC", "", nil, types.WeatherInfo{}, NavCanda, "", 0},

	{"Fond du Lac", "CZFD", Highways + "fonddulac", hc(2), types.WeatherInfo{}, Highways + "/fonddulac", "122.175", 0},
	{"Wollaston", "CZWL", Highways + "wollaston", hc(2), types.WeatherInfo{}, Highways + "/wollaston", "122.075", 0},
	{"Ile A La Crosse", "CJF3", Highways + "ilealacrosse", hc(2), types.WeatherInfo{}, Highways + "/ilealacrosse", "122.075", 0},
	{"Cumberland House", "CJT4", Highways + "cumberlandhouse", hc(2), types.WeatherInfo{}, Highways + "/cumberlandhouse", "122.975", 0},
	{"La Loche", "CJL4", Highways + "laloche", hc(2), types.WeatherInfo{}, Highways + "/laloche", "122.975", 0},
	{"Patuanak", "CKB2", Highways + "patuanak", hc(1), types.WeatherInfo{}, "", "", 0},
	{"Pelican Narrows", "CJW4", Highways + "pelican", hc(1), types.WeatherInfo{}, "", "", 0},
	{"Pinehouse", "CZPO", Highways + "pinehouse", hc(2), types.WeatherInfo{}, Highways + "/pinehouse", "123.175", 0},
	{"Buffalo Narrows", "CYVT", Highways + "buffalonarrows", hc(2), types.WeatherInfo{}, NavCanda, "", 0},
	{"Hudson Bay", "CYHB", Highways + "hudsonbay", hc(2), types.WeatherInfo{}, "", "", 0},
	{"Stony Rapids", "CYSF", Highways + "stonyrapids", hc(2), types.WeatherInfo{}, NavCanda, "", 0},
	{"Sandy Bay", "CJY4", Highways + "sandybay", []string{"/ptz.jpg", "/ptz2.jpg", "/ptz3.jpg"}, types.WeatherInfo{}, Highways + "/sandybay", "122.550", 0},
	{"Meadow Lake", "CYLJ", Highways + "meadowlake", hc(2), types.WeatherInfo{}, NavCanda, "", 0},
	{"Uranium City", "CYBE", Highways + "uranium", hc(1), types.WeatherInfo{}, "https://metar-taf.com/CWDC", "", 0},

	{"Charlot River", "CJP9", "http://saskpower.glmobile.com/charlot", []string{"/runway.jpg", "/hill.jpg"}, types.WeatherInfo{}, "", "", 0},

	{"Flin Flon", "CYFO", MetCam + "CYFO", []string{"/CYFO_SW-full-e.jpeg", "/CYFO_NW-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, "", 0},
	{"North Battleford", "CYQW", MetCam + "CYQW", []string{"/CYQW_S-full-e.jpeg", "/CYQW_W-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, "", 0},

	{"Leismer", "CET2", "", nil, types.WeatherInfo{}, "https://cet2.ca/CET2_metar", "", 0},
	{"Cold Lake", "CYOD", "", nil, types.WeatherInfo{}, NavCanda, "", 0},

	{"Regina", "CYQR", "", nil, types.WeatherInfo{}, NavCanda, "", 0},

	{"Fort Mac", "CYMM", MetCam + "CYMM", []string{"/CYMM_SE-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, "", 0},
	{"Fort Smith", "CYSM", "", nil, types.WeatherInfo{}, NavCanda, "", 0},
	{"Fort Chip", "CYPY", MetCam + "CYPY", []string{"/CYPY_NW-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, "", 0},

	// note: new req for v2: map single source to multiple airports (cold lake and grace lake are the same, CYOD)
	{"The Pas", "CYQD", "", nil, types.WeatherInfo{}, NavCanda, "", 0},
	{"Lloydminster", "CYLL", MetCam + "CYLL", []string{"/CYLL_SE-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, "", 0},

	{"Swift Current", "CYYN", MetCam + "CYYN", []string{"/CYYN_S-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, "", 0},
	{"Medicine Hat", "CYXH", "", nil, types.WeatherInfo{}, NavCanda, "", 0},
	{"Thompson", "CYTH", "", nil, types.WeatherInfo{}, NavCanda, "", 0},
	{"Yorkton", "CYQV", MetCam + "CYQV", []string{"/CYQV_S-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, "", 0},

	{"Lynn Lake", "CYYL", MetCam + "CYYL", []string{"/CYYL_E-full-e.jpeg", "/CYYL_S-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, "", 0},

	{"Seabee Mine", "CCB2", "http://ssrmining.glmobile.com/" + "seabee", []string{"/ptz1.jpg"}, types.WeatherInfo{}, "http://ssrmining.glmobile.com/seabee/", "", 0},

	{"Grace Lake", "CJR3", "", nil, types.WeatherInfo{}, "https://metar-taf.com/CJR3", "", 0},
}
var indexData = types.IndexData{
	AirportInformation: airportInfo,
	GenericPageData: types.GenericPageData{
		LastUpdate: time.Time{},
		Version:    globals.Version,
	},
}

var gfaData = types.GfaPageData{
	GenericPageData: types.GenericPageData{
		LastUpdate: time.Time{},
		Version:    globals.Version,
	},
}

var notamData = types.NotamPageData{
	NoTamOptions: []string{"CJW7", "CKQ8", "CYKC", "CYKJ", "CYNL", "CYXE",
		"CYPA", "CYVC", "CZFD", "CZWL", "CJF3", "CJT4",
		"CJL4", "CKB2", "CJW4", "CZPO", "CYVT", "CYHB",
		"CYSF", "CJY4", "CYLJ", "CYBE", "CJP9", "CYFO", "CYQW",
		"CET2", "CYOD", "CYQR",
		"CYMM", "CYSM", "CYPY", "CYQD", "CYLL",
		"CYYN", "CYXH", "CYTH", "CYQV", "CCB2", "CJR3",
	},
	GenericPageData: types.GenericPageData{
		LastUpdate: time.Time{},
		Version:    globals.Version,
	},
}

var windsData = types.WindsPageData{
	WindsOptions: []string{"CYQR", "CYVC", "CYXE", "CYYL"},
	GenericPageData: types.GenericPageData{
		LastUpdate: time.Time{},
		Version:    globals.Version,
	},
	MaxInt: math.MaxInt,
}

func LoadTemplate(templatePath string, templateName string, templateStr ...string) *template.Template {
	var templateString string

	templateFunctions := template.FuncMap{"formatTaf": func(taf string) template.HTML {
		res := strings.Replace(taf, "\n", "<br>", -1)
		res = strings.Replace(res, " ", "&nbsp;<wbr>", -1)
		return template.HTML(res)
	},
		"sort": func(s []string) []string {
			sorted := make([]string, len(s))
			copy(sorted, s)
			sort.Strings(sorted)
			return sorted
		},
	}

	globals.Logger.Printf("Loading %s Template", templateName)
	if globals.Env == "local" && templatePath != "" {
		globals.Logger.Printf("Getting %s from disk", templatePath)
		file, err := os.ReadFile(templatePath)
		if err != nil {
			globals.Logger.Printf(err.Error())
		}
		templateString = string(file)
	} else if len(templateStr) > 0 {
		globals.Logger.Printf("Loading template from templateStr param")
		return template.Must(template.New(templateName).Funcs(templateFunctions).Parse(templateStr[0]))
	}

	return template.Must(template.New(templateName).Funcs(templateFunctions).Parse(templateString))
}

func TryUpdateGFAData() {
	start := time.Now()
	gfaData.Lock()
	defer gfaData.Unlock()

	if !(time.Since(gfaData.LastUpdate) > 30*time.Second) {
		return
	}

	ids, err := pull.GetGFAImageIds()
	if err != nil {
		gfaData.Error = err
		gfaData.GfaInfo = types.GfaInfo{}
		return
	}

	gfaData.GfaInfo = ids
	gfaData.Error = err
	gfaData.LastUpdate = time.Now().UTC()

	globals.Logger.Printf("Updated GFA data in %d ms", time.Since(start).Milliseconds())
}

// cameco is like 5->50 seconds to update so im just gonna update that every 2 minutes instead of users see errors all the time
func TryCamecoUpdate() {
	start := time.Now()
	dataChan := make(chan types.WeatherPullInfo)
	var wg sync.WaitGroup
	wg.Add(1)
	pull.GetAllCamecoData(dataChan, &wg)

	// close chan so read loop doesn't hang
	go func() {
		wg.Wait()
		close(dataChan)
	}()

	var metarData []types.WeatherPullInfo
	for metar := range dataChan {
		metarData = append(metarData, metar)
	}

	indexData.Lock()
	defer indexData.Unlock()

	// Update current data w/ new pulled data
	for _, pulledAirport := range metarData {
		for j, currentDataAirports := range indexData.AirportInformation {
			if currentDataAirports.AirportCode == pulledAirport.AirportCode {
				indexData.AirportInformation[j].Error = pulledAirport.Error

				if pulledAirport.Error != nil {
					indexData.AirportInformation[j].Stale += 2
					continue
				} else {
					indexData.AirportInformation[j].Stale = 0
				}

				indexData.AirportInformation[j].Metar = pulledAirport.Metar
				indexData.AirportInformation[j].Taf = pulledAirport.Taf
			}
		}
	}

	indexData.CamecoLastUpdated = time.Now().UTC()
	globals.Logger.Printf("Updated Cameco METAR data in %d ms", time.Since(start).Milliseconds())
}

func TryNavCanUpdate() {
	start := time.Now()
	dataChan := make(chan types.WeatherPullInfo)
	var wg sync.WaitGroup
	wg.Add(1)
	go pull.GetNavCanadaMetars(dataChan, &wg)

	// close chan so read loop doesn't hang
	go func() {
		wg.Wait()
		close(dataChan)
	}()

	var metarData []types.WeatherPullInfo
	for metar := range dataChan {
		metarData = append(metarData, metar)
	}

	indexData.Lock()
	defer indexData.Unlock()

	// Update current data w/ new pulled data
	for _, pulledAirport := range metarData {
		for j, currentDataAirports := range indexData.AirportInformation {
			if currentDataAirports.AirportCode == pulledAirport.AirportCode ||
				currentDataAirports.AirportCode == "CJR3" && pulledAirport.AirportCode == "CYOD" {
				indexData.AirportInformation[j].Error = pulledAirport.Error

				if pulledAirport.Error != nil {
					// magic numbers ill fix once i rewrite everything, 2 b/c cycles are every 2min
					indexData.AirportInformation[j].Stale += 2
					continue
				} else {
					indexData.AirportInformation[j].Stale = 0
				}

				indexData.AirportInformation[j].Metar = pulledAirport.Metar
				indexData.AirportInformation[j].Taf = pulledAirport.Taf
			}
		}
	}

	globals.Logger.Printf("Updated NavCan METAR data in %d ms", time.Since(start).Milliseconds())
}

func TryUpdateMETARData() {
	start := time.Now()
	indexData.Lock()
	defer indexData.Unlock()

	if !(time.Since(indexData.LastUpdate) > 30*time.Second) {
		return
	}

	dataChan := make(chan types.WeatherPullInfo)

	var wg sync.WaitGroup
	wg.Add(3)
	go pull.GetAllHighwayData(dataChan, &wg)
	go pull.GetPointsNorthMetar(dataChan, &wg)
	go pull.GetAllMesotech(dataChan, &wg)

	go func(wg *sync.WaitGroup) {
		wg.Add(1)
		pull.GetEnvironmentCanada("CYBE", dataChan)
		wg.Done()
	}(&wg)

	// close chan so read loop doesn't hang
	go func() {
		wg.Wait()
		close(dataChan)
	}()

	var metarData []types.WeatherPullInfo
	for metar := range dataChan {
		metarData = append(metarData, metar)
	}

	indexData.LastUpdate = time.Now().UTC()
	// Update current data w/ new pulled data
	for _, pulledAirport := range metarData {
		for j, currentDataAirports := range indexData.AirportInformation {
			if currentDataAirports.AirportCode == pulledAirport.AirportCode {
				indexData.AirportInformation[j].Metar = pulledAirport.Metar
				indexData.AirportInformation[j].Taf = pulledAirport.Taf
				indexData.AirportInformation[j].Error = pulledAirport.Error
				// only highways needs this
				if pulledAirport.UpdatedImageUrls != nil {
					globals.Logger.Printf("updating image urls for %s to %v", indexData.AirportInformation[j].AirportName, *pulledAirport.UpdatedImageUrls)
					indexData.AirportInformation[j].CamPicUrls = *pulledAirport.UpdatedImageUrls
				}
			}
		}
	}

	globals.Logger.Printf("Updated METAR data in %d ms", time.Since(start).Milliseconds())
}
