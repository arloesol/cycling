package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"encoding/json"
	"net/http"

	"github.com/gocolly/colly"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// based on json response and using https://transform.tools/json-to-go

type RouteList struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	Pagesize int `json:"pagesize"`
	Types    struct {
		Wijn              []interface{} `json:"wijn"`
		Druif             []interface{} `json:"druif"`
		Bloesem           []interface{} `json:"bloesem"`
		Fietsen           []string      `json:"fietsen"`
		BrabantsTrekpaard []interface{} `json:"brabants trekpaard"`
		Hop               []interface{} `json:"hop"`
		Ontdekken         []interface{} `json:"ontdekken"`
		Oorlog            []interface{} `json:"oorlog"`
		Bier              []interface{} `json:"bier"`
		Kinderen          []interface{} `json:"kinderen"`
		NatuurRecreatie   []string      `json:"natuur & recreatie"`
	} `json:"types"`
	Filter struct {
		Accessibility  []interface{} `json:"accessibility"`
		Cities         []interface{} `json:"cities"`
		Distance       []interface{} `json:"distance"`
		Geofilter      interface{}   `json:"geofilter"`
		Page           int           `json:"page"`
		Pagesize       int           `json:"pagesize"`
		Themes         []string      `json:"themes"`
		Types          []interface{} `json:"types"`
		Region         string        `json:"region"`
		DistanceRanges []interface{} `json:"distanceRanges"`
		GeoFilter      interface{}   `json:"geoFilter"`
		PageSize       int           `json:"pageSize"`
	} `json:"filter"`
	Facets struct {
		Facets []struct {
			Name   string `json:"name"`
			Values []struct {
				Value string `json:"value"`
				Times int    `json:"times"`
			} `json:"values"`
		} `json:"facets"`
	} `json:"facets"`
	Results struct {
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
	routelist        RouteList
	routenbr         int
	date             = ""
	name             = ""
	title            = ""
	description      = ""
	length           = ""
	gpxfile          = ""
	nodes            = []string{}
	startpunt        = ""
	content          = ""
	sideimages       = 0
	spaceStartLineRE = regexp.MustCompile("\\n[ \\t]*")
	NewLineRE        = regexp.MustCompile("\\n+")
	NodesRE          = regexp.MustCompile(" |\\.")
	caser            = cases.Title(language.English)
)

func main() {
	c := colly.NewCollector()

	date = time.Now().UTC().Format("2006-01-02")

	os.Mkdir("gpx", 0750)
	os.Mkdir("gpx/vlaamsbrabant", 0750)
	os.Mkdir("img", 0750)
	os.Mkdir("img/gallery", 0750)
	os.Mkdir("img/page", 0750)
	os.Mkdir("route", 0750)
	os.Mkdir("route/vlaamsbrabant", 0750)

	// get primary content and nodes
	c.OnHTML("div.pdintro__content > ul", func(e *colly.HTMLElement) {
		e.ForEach("li", func(nbr int, e *colly.HTMLElement) {
			txtfragment := cleanupTxt(e.Text)
			e.ForEach("a", func(nbr int, e *colly.HTMLElement) {
				link := e.Attr("href")
				linktxt := e.Text
				if linktxt != "" {
					txtfragment = strings.Replace(txtfragment, linktxt, "["+linktxt+"]("+e.Request.AbsoluteURL(link)+")", 1)
				}
			})
			content += txtfragment + "\n\n"
			nodeprefix := "Volg de knooppunten:"
			if strings.HasPrefix(txtfragment, nodeprefix) {
				// vb: Volg de knooppunten: 61 - 25 - 62 - 84 - 85 - 88 - 24 - 23 - 64 - 15 - 21 - 18 - 58 - 69 - 61.
				nodetxt := txtfragment[len(nodeprefix):]
				nodetxt = NodesRE.ReplaceAllString(nodetxt, "") // remove " " and "."
				nodes = strings.Split(nodetxt, "-")
				// fmt.Println("nodes: ", nodes)
			}
		})
		// fmt.Println(content)
	})

	// get gpx file details
	// example: <div class="btnfield left"><a class="vlbr-btn btn-green2 matomo_download" href="https://tools.nodemapp.com/webservices/download-gpx/329" target="_blank">Download de route als GPX</a></div>
	c.OnHTML("div.btnfield > a.matomo_download", func(e *colly.HTMLElement) {
		if e.Text == "Download de route als GPX" {
			gpxurl := e.Attr("href")
			gpxfile = "gpx/vlaamsbrabant/" + title + ".gpx"
			ctx := colly.NewContext()
			ctx.Put("filename", gpxfile)
			// fmt.Println("fetching gpx ", gpxfile, gpxurl)
			c.Request("GET", e.Request.AbsoluteURL(gpxurl), nil, ctx, nil)
		}
	})

	// get starting point
	c.OnHTML("div.pdintro__details > ul.pdintro__details__content", func(e *colly.HTMLElement) {
		e.ForEach("li", func(nbr int, e *colly.HTMLElement) {
			txt := strings.TrimSpace(e.Text)
			if strings.HasPrefix(txt, "Vertrekplaats") {
				startpunt = cleanupTxt(txt[len("Vertrekplaats"):])
				// fmt.Println("startpunt:", startpunt)
			}
		})
	})

	// get side images
	c.OnHTML("div.pdintro__details > ul.pdintro__medialist > li.pdintro__media > img", func(e *colly.HTMLElement) {
		// fmt.Println("got a side image", e.Attr("src"))
		imgurl := e.Attr("src")
		imgfilename := imgname(imgurl, true)
		if imgfilename != "" {
			ctx := colly.NewContext()
			ctx.Put("filename", "img/gallery/"+name+"/"+imgfilename)
			c.Request("GET", e.Request.AbsoluteURL(imgurl), nil, ctx, nil)
		}
	})

	// save gpx and image files
	c.OnResponse(func(r *colly.Response) {
		filename := r.Ctx.Get("filename")
		if filename != "" {
			fmt.Println("saving file", filename)
			r.Save(filename)
		}
	})

	getFietsRoutes()

	for i, route := range routelist.Results.Results {
		routenbr = i
		fmt.Println("visiting webpage for ", route.Title, " at ", route.URL)
		content = ""
		nodes = nil
		gpxfile = ""
		startpunt = ""
		sideimages = 0
		title = strings.ReplaceAll(route.Title, " ", "_")
		name = "be.vlaamsbrabant." + title
		os.Mkdir("img/gallery/"+name, 0750)
		c.Visit("https://www.toerismevlaamsbrabant.be/" + route.URL)
		getMainImage(route.Image, c)
		createRoutePage(i)
	}
}

func getMainImage(url string, c *colly.Collector) {
	ctx := colly.NewContext()
	ctx.Put("filename", "img/gallery/"+name+"/"+imgname(url, false))
	c.Request("GET", "https://www.toerismevlaamsbrabant.be"+url, nil, ctx, nil)
}

func createRoutePage(i int) {
	route := routelist.Results.Results[i]
	//fmt.Println(route)

	f, _ := os.Create("route/vlaamsbrabant/" + name + ".md")
	defer f.Close()
	mdContent := `---
title: "%s"
subtitle: "%s"
date: "%s"
description: "%s" 
tags:
- flanders
- medium
categories: 
- route
- official
region: "flanders"
source: "be.vlaamsbrabant"
ext_url: "%s"
gpx: "vlaamsbrabant/%s"
length: %s%s
---

## Let's Go ! 

%s

## Start

%s
`
	nodestr := ""
	if nodes != nil {
		nodestr = "\nnodetype: \"vlaams\"\nnodes: \"" + strings.Join(nodes, ",") + "\""
	}
	if gpxfile == "" {
		// skip this route - cleanup
		fmt.Println("no gpxfile - no route page - cleaning up")
		os.RemoveAll("img/gallery/" + name)
	} else {
		f.WriteString(fmt.Sprintf(mdContent, title, route.Summary, date, route.Summary, "https://www.toerismevlaamsbrabant.be"+route.URL, title+".gpx", fmt.Sprintf("%.f", route.Distances[0]), nodestr, content, startpunt))
	}
}

func getFietsRoutes() {
	// use the api of the website to get the basic bike-route info in JSON format
	payload := []byte(`{"types":[],"accessibility":[],"cities":[],"themes":["Fietsroutes"],"distance":[],"region":"ALL","page":0,"pagesize":500}`)
	url := "https://www.toerismevlaamsbrabant.be/api/catalogus/251"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		//jsonStr := string(body)
		//fmt.Println("Response: ", jsonStr)
		err = json.Unmarshal(body, &routelist)
		if err != nil {
			panic(err)
		}
		fmt.Println("got back a list of", len(routelist.Results.Results), "routes")
	} else {
		fmt.Println("Post to site failed with error: ", resp.Status)
		os.Exit(1)
	}
}

// /Images/diepensteyn-fietsroute-1_tcm251-121024_w640_n.jpg -> title + "_" + idx
func imgname(url string, side bool) string {
	slice := strings.Split(url, "/")
	filename := slice[len(slice)-1]
	slice = strings.Split(filename, ".")
	if len(slice) == 1 {
		return ""
	}
	//base := strings.Join(slice[:len(slice)-1],".")
	ext := slice[len(slice)-1]
	if side {
		filename = title + "_" + fmt.Sprintf("%d", sideimages) + "." + ext
		sideimages += 1
	} else {
		filename = title + "." + ext
	}

	return filename
}

func cleanupTxt(s string) string {
	s = strings.ReplaceAll(s, "\n", "\n\n")
	s = strings.ReplaceAll(s, " ", " ")
	s = spaceStartLineRE.ReplaceAllString(s, "\n")
	s = NewLineRE.ReplaceAllString(s, "\n\n")
	s = strings.TrimSpace(s)

	return s
}
