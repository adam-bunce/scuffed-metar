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
	"strconv"
	"strings"
	"sync"
	"time"
)

// TODO: semantic compression for all these

//go:embed static
var files embed.FS
var fileServer = http.FileServer(http.FS(files))

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
	{"Cigar Lake", "CJW7", "", nil, types.WeatherInfo{}, Cameco + "CJW7/", ""},
	{"McArthur River", "CKQ8", "", nil, types.WeatherInfo{}, Cameco + "CKQ8/", ""},
	{"Collins Bay / Rabbit Lake", "CYKC", "", nil, types.WeatherInfo{}, Cameco + "CYKC/", ""},
	{"Key Lake", "CYKJ", MetCam + "CYKJ", []string{"/CYKJ_SW-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, ""},

	{"Points North", "CYNL", "", nil, types.WeatherInfo{}, PointsNorth + "CYNL_metar.html", ""},

	{"Saskatoon (Stoon)", "CYXE", "", nil, types.WeatherInfo{}, NavCanda, ""},
	{"Prince Albert (PA)", "CYPA", "", nil, types.WeatherInfo{}, NavCanda, ""},
	{"La Ronge", "CYVC", "", nil, types.WeatherInfo{}, NavCanda, ""},

	{"Fond du Lac", "CZFD", Highways + "fonddulac", hc(2), types.WeatherInfo{}, Highways + "/fonddulac", "122.175"},
	{"Wollaston", "CZWL", Highways + "wollaston", hc(2), types.WeatherInfo{}, Highways + "/wollaston", "122.075"},
	{"Ile A La Crosse", "CJF3", Highways + "ilealacrosse", hc(2), types.WeatherInfo{}, Highways + "/ilealacrosse", "122.075"},
	{"Cumberland House", "CJT4", Highways + "cumberlandhouse", hc(2), types.WeatherInfo{}, Highways + "/cumberlandhouse", "122.975"},
	{"La Loche", "CJL4", Highways + "laloche", hc(2), types.WeatherInfo{}, Highways + "/laloche", "122.975"},
	{"Patuanak", "CKB2", Highways + "patuanak", hc(1), types.WeatherInfo{}, "", ""},
	{"Pelican Narrows", "CJW4", Highways + "pelican", hc(1), types.WeatherInfo{}, "", ""},
	{"Pinehouse", "CZPO", Highways + "pinehouse", hc(2), types.WeatherInfo{}, Highways + "/pinehouse", "123.175"},
	{"Buffalo Narrows", "CYVT", Highways + "buffalonarrows", hc(2), types.WeatherInfo{}, NavCanda, ""},
	{"Hudson Bay", "CYHB", Highways + "hudsonbay", hc(2), types.WeatherInfo{}, "", ""},
	{"Stony Rapids", "CYSF", Highways + "stonyrapids", hc(2), types.WeatherInfo{}, NavCanda, ""},
	{"Sandy Bay", "CJY4", Highways + "sandybay", []string{"/ptz.jpg", "/ptz2.jpg", "/ptz3.jpg"}, types.WeatherInfo{}, Highways + "/sandybay", "122.550"},
	{"Meadow Lake", "CYLJ", Highways + "meadowlake", hc(2), types.WeatherInfo{}, NavCanda, ""},
	{"Uranium City", "CYBE", Highways + "uranium", hc(1), types.WeatherInfo{}, "", ""},

	{"Charlot River", "CJP9", "http://saskpower.glmobile.com/charlot", []string{"/runway.jpg", "/hill.jpg"}, types.WeatherInfo{}, "", ""},

	{"Flin Flon", "CYFO", MetCam + "CYFO", []string{"/CYFO_SW-full-e.jpeg", "/CYFO_NW-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, ""},
	{"North Battleford", "CYQW", MetCam + "CYQW", []string{"/CYQW_S-full-e.jpeg", "/CYQW_W-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, ""},

	{"Leismer", "CET2", "", nil, types.WeatherInfo{}, "https://cet2.ca/CET2_metar", ""},
	{"Christina Lake", "CCL3", "", nil, types.WeatherInfo{}, "https://ccl3.azurewebsites.net/", ""},

	{"Regina", "CYQR", "", nil, types.WeatherInfo{}, NavCanda, ""},

	{"Fort Mac", "CYMM", MetCam + "CYMM", []string{"/CYMM_SE-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, ""},
	{"Fort Smith", "CYSM", "", nil, types.WeatherInfo{}, NavCanda, ""},
	{"Fort Chip", "CYPY", MetCam + "CYPY", []string{"/CYPY_NW-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, ""},

	{"The Pas", "CYQD", "", nil, types.WeatherInfo{}, NavCanda, ""},
	{"Lloydminster", "CYLL", "", nil, types.WeatherInfo{}, NavCanda, ""},

	{"Swift Current", "CYYN", MetCam + "CYYN", []string{"/CYYN_S-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, ""},
	{"Medicine Hat", "CYXH", "", nil, types.WeatherInfo{}, NavCanda, ""},
	{"Thompson", "CYTH", "", nil, types.WeatherInfo{}, NavCanda, ""},
	{"Yorkton", "CYQV", MetCam + "CYQV", []string{"/CYQV_S-full-e.jpeg"}, types.WeatherInfo{}, NavCanda, ""},
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
		"CET2", "CCL3", "CYQR",
		"CYMM", "CYSM", "CYPY", "CYQD", "CYLL",
		"CYYN", "CYXH", "CYTH", "CYQV",
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
				indexData.AirportInformation[j].Metar = pulledAirport.Metar
				indexData.AirportInformation[j].Taf = pulledAirport.Taf
				indexData.AirportInformation[j].Error = pulledAirport.Error
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
			if currentDataAirports.AirportCode == pulledAirport.AirportCode {
				indexData.AirportInformation[j].Metar = pulledAirport.Metar
				indexData.AirportInformation[j].Taf = pulledAirport.Taf
				indexData.AirportInformation[j].Error = pulledAirport.Error
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
