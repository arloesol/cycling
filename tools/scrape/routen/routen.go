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
		Pagetoparse: "https://www.routen.be/langs-gentse-wateren-fietsroute",
		Savegpx:     true,
		Saveimg:     true,
		Source:      "routen",
		Srcpfx:      "be.routen.",
		Tags:        []string{"flanders"},
		Categories:  []string{"official"},
		Region:      "flanders",
		NodeType:    "flanders",
	}
	route lib.Route

	routenbr = 0
)

func main() {
	c := colly.NewCollector()

	lib.Mkalldirs(cfg)

	// route page urls from overview page
	c.OnHTML("div.view-content > div > div > article > div", func(e *colly.HTMLElement) {
		route = lib.Emptyroute
		// short description - Title
		e.ForEach("div.carousel__title > h3", func(nbr int, e *colly.HTMLElement) {
			route.Description = lib.CleanTxt(e.Text)
		})
		// length
		e.ForEach("div.carousel__property.fas.fa-ruler > span.route-lengt", func(nbr int, e *colly.HTMLElement) {
			route.Length = lib.TxttoInt(e.Text)
		})
		// route page url
		e.ForEach("div.carousel__link > a", func(nbr int, e *colly.HTMLElement) {
			route.Routeurl = e.Request.AbsoluteURL(e.Attr("href"))
			lib.Routename(cfg, &route, e.Attr("href")[1:])
			route.Title = lib.NametoTitle(e.Attr("href")[1:])
		})
		routenbr += 1
		if cfg.Pagetoparse == "" || route.Routeurl == cfg.Pagetoparse {
			fmt.Println("route:", route.Routeurl)
			lib.Mkdirs(cfg, route)
			c.Request("GET", route.Routeurl, nil, nil, nil)
		}
	})

	// main content
	c.OnHTML("#route-content > div > div > div > p", func(e *colly.HTMLElement) {
		route.Content = strings.ReplaceAll(lib.CleanTxt(e.Text), "\"", "")
		route.Description = lib.Firstline((route.Content))
	})

	// gpx file
	c.OnHTML("#main-content > div.node.node--type-route.node--view-mode-full.ds-1col.clearfix > div.field.field--name-route-header.field--type-ds.field--label-hidden.field__item > div > div > div.full-map__left > div.full-map__header > div > div > div.field.field--name-route-actions.field--type-ds.field--label-hidden.field__item > a", func(e *colly.HTMLElement) {
		downloadsdir := e.Attr("href")
		downloadsdir = downloadsdir[:len(downloadsdir)-1]
		lib.SaveGPX(c, e, cfg, &route, e.Request.AbsoluteURL(downloadsdir+"/gpx"))
	})

	// route segments/POI
	c.OnHTML("#main-content > div.node.node--type-route.node--view-mode-full.ds-1col.clearfix > div.container.route-description > div > div > article", func(e *colly.HTMLElement) {
		poi := lib.RoutePOI{}
		// main title
		e.ForEach("h3", func(nbr int, e *colly.HTMLElement) {
			poi.Title = e.Text
		})
		// subtitle
		e.ForEach("em", func(nbr int, e *colly.HTMLElement) {
			poi.Content = strings.TrimSpace(e.Text) + "\n\n"
		})
		// POI/segment description
		e.ForEach("section > div.route-segment__description > p", func(nbr int, e *colly.HTMLElement) {
			poi.Content += strings.TrimSpace(e.Text)
		})
		// image
		e.ForEach("section > div.route-segment__media > div > div", func(nbr int, e *colly.HTMLElement) {
			poi.Imgurl = strings.Split(e.Attr("data-flickity-bg-lazyload"), "?")[0]
		})
		route.POIs = append(route.POIs, poi)
	})

	// save gpx file
	c.OnResponse(lib.SaveOnResponse(cfg))

	// route page scrape -> create route .md page
	c.OnScraped(func(r *colly.Response) {
		if route.Gpxfile != "" {
			lib.Routepage(cfg, route)
		}
	})

	// visit all route ovv pages until there are no more (routes)
	oldnbrroutes := -1
	for page := 0; routenbr > oldnbrroutes; page++ {
		oldnbrroutes = routenbr
		ovvurl := "https://www.routen.be/fietsroutes?page=" + fmt.Sprintf("%d", page)
		fmt.Println(ovvurl)
		c.Visit(ovvurl)
	}

}
