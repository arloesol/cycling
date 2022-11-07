package main

import (
	"fmt"
	"regexp"
	"strings"

	"cycling.io/m/v2/lib"
	"github.com/gocolly/colly"
)

// based on json response and using https://transform.tools/json-to-go -- removed sections not needed

type RouteList struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	Pagesize int `json:"pagesize"`
	Results  struct {
		Results []struct {
			ID            string        `json:"id"`
			URL           string        `json:"url"`
			Image         string        `json:"image"`
			Summary       string        `json:"summary"`
			Title         string        `json:"title"`
			Latitude      float64       `json:"latitude"`
			Longitude     float64       `json:"longitude"`
			Locations     []string      `json:"locations"`
			Accessibility []interface{} `json:"accessibility"`
			Distances     []float64     `json:"distances"`
			Cities        []string      `json:"cities"`
		} `json:"results"`
	} `json:"results"`
}

var (
	cfg = lib.Cfg{
		Pagetoparse: "https://www.toerismevlaamsbrabant.be/producten/fietsen/fietsproducten/diepensteyn-fietsroute/index.html",
		Savegpx:     true,
		Saveimg:     true,
		Source:      "vlaamsbrabant",
		Srcpfx:      "be.vlaamsbrabant.",
	}
	route lib.Route

	routelist RouteList
	NodesRE   = regexp.MustCompile(` |\.`)
)

func main() {
	c := colly.NewCollector()

	lib.Mkalldirs(cfg)

	// get primary content and nodes
	c.OnHTML("div.pdintro__content > ul", func(e *colly.HTMLElement) {
		e.ForEach("li", func(nbr int, e *colly.HTMLElement) {
			txtfragment := lib.Txtandlinks(e)
			route.Content += txtfragment + "\n\n"
			nodeprefix := "Volg de knooppunten:"
			if strings.HasPrefix(txtfragment, nodeprefix) {
				// vb: Volg de knooppunten: 61 - 25 - 62 - 84 - 85 - 88 - 24 - 23 - 64 - 15 - 21 - 18 - 58 - 69 - 61.
				nodetxt := txtfragment[len(nodeprefix):]
				nodetxt = NodesRE.ReplaceAllString(nodetxt, "") // remove " " and "."
				route.Nodes = strings.Split(nodetxt, "-")
				// fmt.Println("nodes: ", nodes)
			}
		})
		// fmt.Println(content)
	})

	// get gpx file details
	c.OnHTML("div.btnfield > a.matomo_download", func(e *colly.HTMLElement) {
		if e.Text == "Download de route als GPX" {
			lib.SaveGPX(c, e, cfg, &route)
		}
	})

	// get starting point
	c.OnHTML("div.pdintro__details > ul.pdintro__details__content", func(e *colly.HTMLElement) {
		e.ForEach("li", func(nbr int, e *colly.HTMLElement) {
			txt := strings.TrimSpace(e.Text)
			if strings.HasPrefix(txt, "Vertrekplaats") {
				route.Startpunt = lib.CleanTxt(txt[len("Vertrekplaats"):])
			}
		})
	})

	// get side images
	c.OnHTML("div.pdintro__details > ul.pdintro__medialist > li.pdintro__media > img", func(e *colly.HTMLElement) {
		// fmt.Println("got a side image", e.Attr("src"))
		lib.SaveIMGanchor(c, e, cfg, &route, '1')
	})

	c.OnResponse(lib.SaveOnResponse)

	getFietsRoutes()

	for _, info := range routelist.Results.Results {
		fmt.Println("visiting webpage for ", info.Title, " at ", info.URL)
		route = lib.Emptyroute
		route.Routeurl = "https://www.toerismevlaamsbrabant.be" + info.URL
		if cfg.Pagetoparse == "" || cfg.Pagetoparse == route.Routeurl {
			route.Title = strings.ReplaceAll(info.Title, " ", "_")
			route.Name = cfg.Srcpfx + route.Title
			route.Length = fmt.Sprintf("%.f", info.Distances[0])
			route.Description = info.Summary
			lib.Mkdirs(cfg, route)
			c.Visit(route.Routeurl)
			getMainImage(info.Image, c)
			lib.Routepage(cfg, route)
		}
	}
}

func getMainImage(url string, c *colly.Collector) {
	lib.SaveIMG(c, "https://www.toerismevlaamsbrabant.be"+url, cfg, &route, '0')
}

func getFietsRoutes() {
	// use the api of the website to get the basic bike-route info in JSON format
	payload := []byte(`{"types":[],"accessibility":[],"cities":[],"themes":["Fietsroutes"],"distance":[],"region":"ALL","page":0,"pagesize":500}`)
	url := "https://www.toerismevlaamsbrabant.be/api/catalogus/251"

	err := lib.Restpost(url, payload, &routelist)
	if err != nil {
		panic(err)
	}
}
