package main

// Usage: go run addelevation.go gpxfile.gpx
// Result: in case some elevations are missing creates a new file gpxfile.gpx.withelevation

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"cycling.io/tools/v2/lib"
	"github.com/tkrajina/gpxgo/gpx"
)

var (
	elevationData map[string]float64
)

// example get https://dhm.agiv.be/api/elevation/v1/DHMVMIXED?SrsIn=4326&Locations=2.7272845900998086,51.149848344629255|2.7305674003238605,51.147908976328452|2.9100556998556946,50.935151476429354

func main() {
	zeroes := 0
	allpointshaveelevation := true
	elevationData = make(map[string]float64)

	// command line arguments
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Please provide a gpx file")
		os.Exit(1)
	}

	// read markdown file - filename = first argument of command
	gpxfilepath := args[0]
	fmt.Println("reading file", gpxfilepath)
	g, err := gpx.ParseFile(gpxfilepath)
	lib.IfErrFatal(fmt.Sprintf("can't open gpx file %s", gpxfilepath), err)

	// get all points that require an elevation
	for _, track := range g.Tracks {
		for _, segment := range track.Segments {
			for _, point := range segment.Points {
				x := point.GetElevation()
				elevationData[fmt.Sprintf("%f,%f", point.Longitude, point.Latitude)] = 0
				if x.Null() { // we need to add elevation to the gpx file
					allpointshaveelevation = false
				}
				if x.Value() == 0 {
					zeroes += 1
				}
			}
		}
	}

	if allpointshaveelevation && zeroes < 5 { // don't change gpx file if no elevation to add
		return
	}

	getAllElevations() // get elevations from the website API

	// set elevation values in gpx structure
	for t := range g.Tracks {
		for s := range g.Tracks[t].Segments {
			for p := range g.Tracks[t].Segments[s].Points {
				g.Tracks[t].Segments[s].Points[p].Elevation.SetValue(elevationData[fmt.Sprintf("%f,%f", g.Tracks[t].Segments[s].Points[p].Longitude, g.Tracks[t].Segments[s].Points[p].Latitude)])
			}
		}
	}

	xml, err := gpx.ToXml(g, gpx.ToXmlParams{Version: "1.1", Indent: true})
	lib.IfErrFatal("can't create xml from gpx data struct", err)

	os.WriteFile(gpxfilepath+".withelevation", xml, 0644)
}

func getAllElevations() {
	batch := 100 // don't remove this batching
	// the dhm.agiv.be website gives wrong elevation values back when asking too many points at the same time
	// not sure what the max value is but 100 seems OK
	// sorting of keys is not really needed but nice to have to be able to compare 2 runs of this tool
	keys := make([]string, 0, len(elevationData))
	for k := range elevationData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for len(keys) > batch {
		getDataFromAgiv(keys[:batch])
		keys = keys[batch:]
	}
	if len(keys) > 1 {
		getDataFromAgiv(keys)
	}
}

func getDataFromAgiv(keys []string) {
	var reqdata = `{ "SrsIn": 4326, "LineString": { "coordinates": [ [`
	var values [][]float64

	reqdata += strings.Join(keys, "], [")
	reqdata += `] ], "type": "LineString" } }`

	//fmt.Println(reqdata)

	data := strings.NewReader(reqdata)

	r, err := http.Post("https://dhm.agiv.be/api/elevation/v1/DHMV2/request", "application/json", data)
	lib.IfErrFatal("Can't request elevation from dhm.agiv.be", err)

	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)

	if r.StatusCode != 200 {
		lib.IfErrFatal("", fmt.Errorf("couldn't get elevationdata from agiv.be website %s", string(body)))
	}

	lib.IfErrFatal("Can't extract data from agiv.be response", json.Unmarshal(body, &values))

	for i, data := range values {
		keylat, _ := strconv.ParseFloat(strings.Split(keys[i], ",")[0], 64)
		keylon, _ := strconv.ParseFloat(strings.Split(keys[i], ",")[1], 64)
		if (math.Abs(keylat-data[1]) > 0.00001) || (math.Abs(keylon-data[2]) > 0.00001) {
			fmt.Println(keys[i], data[1], data[2])
		}
		elevationData[keys[i]] = data[3]
	}

	//fmt.Println(string(body))
	//fmt.Println(elevationData)
}
