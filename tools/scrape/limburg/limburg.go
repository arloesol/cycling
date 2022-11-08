package main

import (
	"fmt"
	"strings"

	"cycling.io/m/v2/lib"
	"github.com/gocolly/colly"
)

var (
	cfg = lib.Cfg{
		Pagetoparse: "",
		//Pagetoparse: "https://www.visitlimburg.be/en/route/cross-border-cycling-along-meuse-river",
		Savegpx:    true,
		Saveimg:    true,
		Source:     "limburg",
		Srcpfx:     "be.visitlimburg.",
		Tags:       []string{"flanders"},
		Categories: []string{"official"},
		Region:     "flanders",
		NodeType:   "flanders",
	}
	route lib.Route

	routenbr = 0
)

func main() {
	c := colly.NewCollector()

	lib.Mkalldirs(cfg)

	// route page urls from overview page
	c.OnHTML("div.overview-search-results", func(e *colly.HTMLElement) {
		e.ForEach("a", func(nbr int, e *colly.HTMLElement) {
			route = lib.Emptyroute
			route.Routeurl = e.Request.AbsoluteURL(e.Attr("href"))
			if strings.HasPrefix(route.Routeurl, "https://www.visitlimburg.be/en/route") {
				lib.Routename(cfg, &route, lib.URLend(route.Routeurl))
				routenbr += 1
				if cfg.Pagetoparse == "" || cfg.Pagetoparse == route.Routeurl {
					lib.Mkdirs(cfg, route)
					lib.LogInfo.Println("Visiting", route.Routeurl)
					e.Request.Visit(route.Routeurl)
				}
			}
		})
	})

	// route description
	c.OnHTML("#block-sassy-content > article > header > div > div > div > div:nth-child(2) > div > div > div > div.slide-content-left.col-md-5.col-md-offset-2 > h1", func(e *colly.HTMLElement) {
		route.Title = strings.ReplaceAll(lib.CleanTxt(e.Text), "\"", "")
	})
	// route long description
	c.OnHTML("#block-sassy-content > article > div > div > div.left > div.route--description", func(e *colly.HTMLElement) {
		route.Content = lib.CleanTxt(e.Text)
		route.Description = lib.Firstline(route.Content)
	})

	// route length
	c.OnHTML("#block-sassy-content > article > div > div > div.right > div.route--info > div.route--info-card > ul > li:nth-child(1) > span:nth-child(2)", func(e *colly.HTMLElement) {
		route.Length = lib.TxttoInt(e.Text)
	})

	// startpunt
	c.OnHTML("#block-sassy-content > article > div > div > div.right > div.route--info > div.route--info-card > ul > li:nth-child(3) > span:nth-child(2)", func(e *colly.HTMLElement) {
		route.Startpunt = lib.CleanTxt(e.Text)
	})

	// gpx file
	c.OnHTML("#block-sassy-content > article > div > div > div.right > div.route--info > div.route--info-download > ul > li:nth-child(1) > a", func(e *colly.HTMLElement) {
		lib.SaveGPX(c, e, cfg, &route, "Attr:href")
	})

	// main image
	c.OnHTML("#block-sassy-content > article > header > div > div > div > div:nth-child(1) > div > div > img", func(e *colly.HTMLElement) {
		lib.SaveIMGanchor(c, e, cfg, &route, "src", "Attr-src")
	})

	// on route POIs
	c.OnHTML("#block-sassy-content > article > footer > div.field--name-route-poi > section > section > div > a", func(e *colly.HTMLElement) {
		poi := lib.RoutePOI{}
		poi.Extlink = e.Request.AbsoluteURL(e.Attr("href"))
		e.ForEach("div.accomodation-image > div > img", func(nbr int, e *colly.HTMLElement) {
			poi.ImgALt = e.Attr("alt")
			poi.Imgurl = strings.Split(e.Attr("data-src"), "?")[0]
		})
		e.ForEach("div.accomodation-title", func(nbr int, e *colly.HTMLElement) {
			poi.Title = lib.CleanTxt(e.Text)
		})
		e.ForEach("div.accomodation-description", func(nbr int, e *colly.HTMLElement) {
			poi.Content = lib.CleanTxt(e.Text)
		})
		route.POIs = append(route.POIs, poi)
	})

	// save gpx or image files
	c.OnResponse(lib.SaveOnResponse(cfg))

	// route page scrape -> create route .md page
	c.OnScraped(func(r *colly.Response) {
		if strings.HasPrefix(r.Request.URL.String(), "https://www.visitlimburg.be/en/route/") {
			lib.Routepage(cfg, route)
		}
	})

	// visit all route ovv pages until there are no more (routes)
	oldnbrroutes := -1
	for page := 0; routenbr > oldnbrroutes; page++ {
		oldnbrroutes = routenbr
		ovvurl := "https://www.visitlimburg.be/en/find-cycling-routes?page=" + fmt.Sprintf("%d", page)
		c.Visit(ovvurl)
	}
}
