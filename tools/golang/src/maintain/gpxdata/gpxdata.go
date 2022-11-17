package main

// Usage: go run gpxdata.go route.md
// Result: create route.md.new with new md strcuture and route data based on linked gpx file(s)

import (
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"

	"cycling.io/tools/v2/lib"
	"github.com/tkrajina/gpxgo/gpx"
	"gopkg.in/yaml.v3"
)

const (
	YAML     = 1
	MarkDown = 2
)

// created using https://yaml2go.prasadg.dev/
type RoutePage struct {
	Date        string   `yaml:"date"`
	Title       string   `yaml:"title"`
	Subtitle    string   `yaml:"subtitle"`
	Description string   `yaml:"description"`
	Categories  []string `yaml:"categories"`
	Tags        []string `yaml:"tags"`
	Region      string   `yaml:"region"`
	Website     string   `yaml:"website"`
	ExtUrl      string   `yaml:"ext_url"`
	Routes      []Route  `yaml:"routes"`
}

// Routes
type Route struct {
	Name        string     `yaml:"name"`
	Nodetype    string     `yaml:"nodetype,omitempty"`
	Nodes       []string   `yaml:"nodes,omitempty"`
	Gpx         string     `yaml:"gpx"`
	Length      float64    `yaml:"length"`
	Up          float64    `yaml:"up"`
	Down        float64    `yaml:"down"`
	MinHeight   float64    `yaml:"minheight"`
	MaxHeight   float64    `yaml:"maxheight"`
	MinSlope    float64    `yaml:"minslope"`
	MaxSlope    float64    `yaml:"maxslope"`
	AvgPosSlope float64    `yaml:"avgposslope"`
	AvgNegSlope float64    `yaml:"avgnegslope"`
	SlopeHisto  [5]float64 `yaml:"slopehisto"`
	Effort      float64    `yaml:"effortlevel"`
	MinLat      float64    `yaml:"minlat"`
	MinLon      float64    `yaml:"minlon"`
	MaxLat      float64    `yaml:"maxlat"`
	MaxLon      float64    `yaml:"maxlon"`
}

// slopeinfo
type SlopeT struct {
	length        float64
	elevationdiff float64
	slope         float64
}

func main() {
	// command line arguments
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Please provide a route markdown file")
		os.Exit(1)
	}

	// read markdown file - filename = first argument of command
	mdFileArg := args[0]
	fmt.Println("Reading", mdFileArg)
	mdall, err := os.ReadFile(mdFileArg)
	lib.IfErrFatal(fmt.Sprintf("error opening markdown file %s", mdFileArg), err)

	slices := bytes.Split(mdall, []byte("---\n"))
	if len(slices) != 3 {
		lib.LogError.Fatalln("can't split YAML and text content from markdown file")
	}
	// now: slices[1] = YAML txt and slices[2] is markdown content txt

	// decode YAML data in markdown into routedate struct
	routedata := RoutePage{}
	err = yaml.Unmarshal([]byte(slices[YAML]), &routedata)
	lib.IfErrFatal(fmt.Sprintf("error getting YAML data from %s", mdFileArg), err)

	newroutedata := routedata
	newroutedata.Routes = []Route{}

	for _, route := range routedata.Routes {
		r := Route{Name: route.Name, Gpx: route.Gpx}
		readGPXFile(&r)
		r.Nodetype = route.Nodetype
		r.Nodes = route.Nodes

		// show main info
		//fmt.Println("route:", r.Name)
		//fmt.Printf("bounds:%.4f %.4f %.4f %.4f\n", r.MinLat, r.MinLon, r.MaxLat, r.MaxLon)
		//fmt.Printf("len:%.2fkm  up:%.0fm dn:%.0fm  height: min:%.1fm max:%.1fm\n", r.Length/1000, r.Up, r.Down, r.MinHeight, r.MaxHeight)
		//fmt.Printf("slope: min:%.1f%% max:%.1f%%\n", r.MinSlope, r.MaxSlope)
		//fmt.Printf("overall efort level %.2f\n", r.Effort)
		//fmt.Printf("avg slope: pos:%.2f neg:%.2f\n", r.AvgPosSlope, r.AvgNegSlope)
		//fmt.Println(r.SlopeHisto)
		//fmt.Println()
		newroutedata.Routes = append(newroutedata.Routes, r)
	}

	createRouteFile(mdFileArg+".new", newroutedata, slices[MarkDown])
}

/* basic rules -> slope < -0.01 --> no effort
    every 0.01 slope above this (-0.01) extra effort multiplier ->
	--> 0.04 -> multiplier = * 5 --> multiplier = (slope + 0.01) * 100

	Todo: add something to add some Effort in case long stretches of high effort without recoup streches

	base effort level = flat run of 30 km.
*/
func calcEffort(slopes []SlopeT) float64 {
	var effort float64 = 0

	for _, slope := range slopes {
		if slope.slope > -0.01 {
			effort += slope.length * (slope.slope + 0.01) * 100
		}
	}
	effort = effort / 30000
	return effort
}

func calcSlopes(g *gpx.GPX) (slopes []SlopeT, minslope float64, maxslope float64, avgposslope float64, avgnegslope float64, slopehisto [5]float64) {
	mindist := 50.0
	totdisttimesposslope := 0.0
	totlengthposslope := 0.0
	totdisttimesnegslope := 0.0
	totlengthnegslope := 0.0
	slope := SlopeT{}
	for _, track := range g.Tracks {
		for _, segment := range track.Segments {
			for pointNo, point := range segment.Points {
				if pointNo == 0 {
					continue
				}
				elevantiondiff := point.Elevation.Value() - segment.Points[pointNo-1].Elevation.Value()
				slope.elevationdiff += elevantiondiff
				slope.length += point.Distance2D(&segment.Points[pointNo-1])
				if (slope.length > mindist) || (pointNo == len(segment.Points)-1) {
					if slope.length == 0 {
						slope.slope = 0
					} else {
						slope.slope = slope.elevationdiff / slope.length
					}
					minslope = math.Min(slope.slope, minslope)
					maxslope = math.Max(maxslope, slope.slope)
					if slope.slope > 0 {
						totdisttimesposslope += slope.length * elevantiondiff
						totlengthposslope += slope.length
					}
					if slope.slope < 0 {
						totdisttimesnegslope += slope.length * elevantiondiff
						totlengthnegslope += slope.length
					}
					if slope.slope >= 0.01 { // add to slopehisto
						meanslope := math.RoundToEven(slope.slope * 100)
						index := int(math.Max(math.Min((meanslope-2)/2, 4), 0))
						slopehisto[index] += slope.length
					}
					slopes = append(slopes, slope)
					slope = SlopeT{}
				}
			}
		}
	}
	minslope *= 100
	maxslope *= 100
	avgposslope = totdisttimesposslope / totlengthposslope
	avgnegslope = totdisttimesnegslope / totlengthnegslope
	return
}

func readGPXFile(r *Route) {
	gpxfilepath := os.Getenv("GITDIR") + "/static/gpxfiles/" + r.Gpx
	fmt.Println("reading file", gpxfilepath)
	g, err := gpx.ParseFile(gpxfilepath)
	lib.IfErrFatal(fmt.Sprintf("can't open gpx file %s", gpxfilepath), err)

	bounds := g.Bounds()
	updn := g.UphillDownhill()
	elbounds := g.ElevationBounds()
	r.MinLat = math.Round(bounds.MinLatitude*1000000) / 1000000
	r.MaxLat = math.Round(bounds.MaxLatitude*1000000) / 1000000
	r.MinLon = math.Round(bounds.MinLongitude*1000000) / 1000000
	r.MaxLon = math.Round(bounds.MaxLongitude*1000000) / 1000000
	r.Length = math.Round(g.Length2D())
	r.Up = math.Round(updn.Uphill)
	r.Down = math.Round(updn.Downhill)
	r.MinHeight = math.Round(elbounds.MinElevation)
	r.MaxHeight = math.Round(elbounds.MaxElevation)
	var slopes []SlopeT
	slopes, r.MinSlope, r.MaxSlope, r.AvgPosSlope, r.AvgNegSlope, r.SlopeHisto = calcSlopes(g)
	r.MinSlope = math.Round(r.MinSlope*10) / 10
	r.MaxSlope = math.Round(r.MaxSlope*10) / 10
	r.AvgPosSlope = math.Round(r.AvgPosSlope*10) / 10
	r.AvgNegSlope = math.Round(r.AvgNegSlope*10) / 10
	for i := range r.SlopeHisto {
		r.SlopeHisto[i] = math.Round(r.SlopeHisto[i])
	}
	r.Effort = math.Round(calcEffort(slopes)*1000) / 1000
}

func createRouteFile(file string, r RoutePage, md []byte) {
	var data []byte
	yamlline := []byte("---\n")

	// yaml seperator
	data = yamlline

	// yaml text
	route, err := yaml.Marshal(r)
	lib.IfErrFatal("Can't yaml marshall route data", err)
	data = append(data, route...)

	// yaml seperator
	data = append(data, yamlline...)

	// markdown text
	data = append(data, md...)

	// write in file
	err = os.WriteFile(file, data, 0644)
	lib.IfErrFatal("Can't write to new route.md file", err)
}
