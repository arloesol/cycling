package main

import (
	"fmt"
	"strings"

	"cycling.io/m/v2/lib"
	"github.com/gocolly/colly"
)

var (
	cfg = lib.Cfg{
		//Pagetoparse: "",
		Pagetoparse: "https://www.flandersbybike.com/coastal-route",
		Savegpx:     true,
		Saveimg:     true,
		Source:      "flandersbybike",
		Srcpfx:      "com.flandersbybike.",
		Tags:        []string{"flanders"},
		Categories:  []string{"official"},
		Region:      "flanders",
	}
	route lib.Route
)

func main() {
	c := colly.NewCollector()

	lib.Mkalldirs(cfg)

	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))

		if e.Request.Depth == 1 && strings.HasPrefix(e.Text, "Look at ") {
			route = lib.Emptyroute
			lib.Routename(cfg, &route, lib.URLend(link))
			route.Routeurl = link
			route.Title = lib.CamelCase.String(strings.ReplaceAll(route.Shortname, "-", " "))
			if cfg.Pagetoparse == "" || cfg.Pagetoparse == link {
				lib.Mkdirs(cfg, route)
				fmt.Println("route", route.Name, route.Routeurl)
				e.Request.Visit(link) // check the route page
			}
		}

		if e.Request.Depth == 2 {
			if strings.Contains(e.Text, "GPX") {
				lib.SaveGPX(c, e, cfg, &route)
			}
		}
	})

	c.OnHTML("meta", func(e *colly.HTMLElement) {
		metaname := e.Attr("name")
		if metaname == "description" {
			route.Description = e.Attr("content")
		}
	})

	c.OnHTML("span", func(e *colly.HTMLElement) {
		if e.Request.Depth == 2 && e.Attr("class") == "field-distance__value" {
			route.Length = lib.TxttoInt(e.Text)
		}
	})

	c.OnHTML("img", func(e *colly.HTMLElement) {
		if e.Request.Depth == 2 && strings.Contains(e.Attr("src"), "images") {
			lib.SaveIMGanchor(c, e, cfg, &route, "src", "Attr-alt")
		}
	})

	c.OnHTML("p > span > span > span > span > span > span > span", func(e *colly.HTMLElement) {
		if e.Request.Depth == 2 {
			route.Content = e.Text
		}
	})

	c.OnResponse(lib.SaveOnResponse(cfg))

	// OnScraped is called after all OnHTMLs for a webpage have been processed - if level = 2 -> create a md file with our collected info
	c.OnScraped(func(r *colly.Response) {
		if r.Request.Depth == 2 {
			lib.Routepage(cfg, route)
		}
	})

	url := "https://www.flandersbybike.com/#routes"
	c.Visit(url)
}
