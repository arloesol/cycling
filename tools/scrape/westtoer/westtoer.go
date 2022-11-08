package main

import (
	"strings"

	"cycling.io/m/v2/lib"
	"github.com/gocolly/colly"
)

var (
	cfg = lib.Cfg{
		//Pagetoparse: "",
		Pagetoparse: "https://www.westtoer.be/nl/vivelevelo-rood",
		Savegpx:     true,
		Saveimg:     true,
		Source:      "westtoer",
		Srcpfx:      "be.westtoer.",
		Tags:        []string{"flanders"},
		Categories:  []string{"official"},
		Region:      "flanders",
		NodeType:    "flanders",
	}
	route lib.Route
)

func main() {
	c := colly.NewCollector()

	lib.Mkalldirs(cfg)

	// ovv of routes
	c.OnHTML("a.node--route", func(e *colly.HTMLElement) {
		depth := e.Request.Depth
		if depth == 1 {
			route = lib.Emptyroute
			route.Routeurl = e.Request.AbsoluteURL(e.Attr("href"))
			lib.Routename(cfg, &route, lib.URLend(route.Routeurl))
			if cfg.Pagetoparse == "" || cfg.Pagetoparse == route.Routeurl {
				lib.Mkdirs(cfg, route)
				e.Request.Visit(route.Routeurl)
			}
		}
	})

	// main image
	c.OnHTML("a", func(e *colly.HTMLElement) {
		if strings.HasPrefix(e.Attr("rel"), "lightbox[field_images]") {
			lib.SaveIMGanchor(c, e, cfg, &route, "href", "Attr-href")
		}
	})

	// title
	c.OnHTML("meta[property=\"og:title\"]", func(e *colly.HTMLElement) {
		route.Title = e.Attr("content")
	})

	// distance
	c.OnHTML("div.field--name-field-route-distance", func(e *colly.HTMLElement) {
		e.ForEach("div.even", func(nbr int, e *colly.HTMLElement) {
			route.Length = lib.TxttoInt(e.Text)
		})
	})

	// main route content
	c.OnHTML("div[itemprop=\"description\"]", func(e *colly.HTMLElement) {
		route.Content = lib.CleanTxt(e.Text)
		e.ForEach("a", func(nbr int, e *colly.HTMLElement) {
			route.ContLinks = append(route.ContLinks, [2]string{e.Text, e.Attr("href")})
		})
		// image in description and text in description "Volg deze bordjes:" -> signage
		e.ForEach("img", func(nbr int, e *colly.HTMLElement) {
			if strings.Contains(route.Content, "Volg deze bordjes:") {
				route.Content = strings.Replace(route.Content, "Volg deze bordjes:", "", 1)
				lib.SaveIMGanchor(c, e, cfg, &route, "src", "Attr-src")
			}
		})
	})

	// startpunt
	c.OnHTML("div.field--name-field-starting-point", func(e *colly.HTMLElement) {
		e.ForEach("div.even", func(nbr int, e *colly.HTMLElement) {
			route.Startpunt = lib.CleanTxt(e.Text)
		})
	})

	// Signage
	c.OnHTML("div.field--name-field-junctions", func(e *colly.HTMLElement) {
		e.ForEach("div.even", func(nbr int, e *colly.HTMLElement) {
			route.Signage = lib.CleanTxt(e.Text)
		})
		lib.TxttoNodes(route.Signage, &route)
	})

	// POI on route
	c.OnHTML("div.pane-node-field-on-your-route", func(e *colly.HTMLElement) {
		e.ForEach("div.content", func(nbr int, e *colly.HTMLElement) { // a POI on route
			poi := lib.RoutePOI{}
			e.ForEach("h3", func(nbr int, e *colly.HTMLElement) { // title of on route POI
				poi.Title = lib.CamelCase.String(e.Text)
			})
			e.ForEach("img[typeof=\"foaf:Image\"]", func(nbr int, e *colly.HTMLElement) {
				poi.Imgurl = e.Request.AbsoluteURL(e.Attr("src"))
			})
			e.ForEach("div.field__item > p", func(nbr int, e *colly.HTMLElement) {
				poitextfragment := strings.TrimSpace(e.Text)
				e.ForEach("a", func(nbr int, e *colly.HTMLElement) {
					poi.ContLinks = append(poi.ContLinks, [2]string{e.Text, e.Attr("href")})
				})
				poi.Content += poitextfragment + "\n"
			})
			route.POIs = append(route.POIs, poi)
		})
	})

	// gpx file
	c.OnHTML("div.field--name-field-gpx", func(e *colly.HTMLElement) {
		e.ForEach("div.field__item", func(nbr int, e *colly.HTMLElement) {
			e.ForEach("a", func(nbr int, e *colly.HTMLElement) {
				lib.SaveGPX(c, e, cfg, &route, "Attr-href")
			})
		})
	})

	c.OnResponse(lib.SaveOnResponse(cfg))

	// OnScraped is called after all OnHTMLs for a webpage have been processed - if level = 2 -> create a md file with our collected info
	c.OnScraped(func(r *colly.Response) {
		if r.Request.Depth == 2 {
			route.Description = lib.Firstline(route.Content)
			lib.Routepage(cfg, route)
		}
	})

	url := "https://www.westtoer.be/nl/doen/fietsroutes?field_geofield_latlon_op=10&field_geofield_latlon=&sort_by=title2&items_per_page=900&map=false&nolocation=TRUE"
	c.Visit(url)
}
