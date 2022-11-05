package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var (
	routenbr         = 0
	date             = ""
	name             = ""
	title            = ""
	description      = ""
	longdescription  = ""
	length           = ""
	gpxfile          = ""
	content          = ""
	routeurl         = ""
	spaceStartLineRE = regexp.MustCompile("\\n[ \\t]*")
	NewLineRE        = regexp.MustCompile("\\n+")
	poititle         = ""
	poisubtitle      = ""
	poidesc          = ""
	poiimgurl        = ""
)

func main() {
	c := colly.NewCollector()

	date = time.Now().UTC().Format("2006-01-02")

	os.Mkdir("gpx", 0750)
	os.Mkdir("gpx/routen", 0750)
	os.Mkdir("route", 0750)
	os.Mkdir("route/routen", 0750)
	os.Mkdir("img", 0750)

	// route page urls from overview page
	c.OnHTML("div.view-content > div > div > article > div", func(e *colly.HTMLElement) {
		// short description - Title
		e.ForEach("div.carousel__title > h3", func(nbr int, e *colly.HTMLElement) {
			description = cleanupTxt(e.Text)
		})
		// length
		e.ForEach("div.carousel__property.fas.fa-ruler > span.route-lengt", func(nbr int, e *colly.HTMLElement) {
			length = strings.TrimSpace(strings.Split(e.Text, ",")[0])
		})
		// route page url
		e.ForEach("div.carousel__link > a", func(nbr int, e *colly.HTMLElement) {
			routeurl = e.Request.AbsoluteURL(e.Attr("href"))
			title = e.Attr("href")[1:]
		})
		fmt.Println("route:", title, description, length, "km", routeurl)
		name = "be.routen." + title
		routenbr += 1
		content = ""
		gpxfile = ""
		os.Mkdir("img/gallery/"+name, 0750)
		c.Request("GET", routeurl, nil, nil, nil)
	})

	//longdescription
	c.OnHTML("#route-content > div > div > div > p", func(e *colly.HTMLElement) {
		longdescription = cleanupTxt(e.Text)
		longdescription = strings.ReplaceAll(longdescription, "\"", "")
		fmt.Println(longdescription)
	})

	// gpx file
	c.OnHTML("#main-content > div.node.node--type-route.node--view-mode-full.ds-1col.clearfix > div.field.field--name-route-header.field--type-ds.field--label-hidden.field__item > div > div > div.full-map__left > div.full-map__header > div > div > div.field.field--name-route-actions.field--type-ds.field--label-hidden.field__item > a", func(e *colly.HTMLElement) {
		downloadsdir := e.Attr("href")
		downloadsdir = downloadsdir[:len(downloadsdir)-1]
		gpxurl := e.Request.AbsoluteURL(downloadsdir + "/gpx")
		gpxfile = title + ".gpx"
		ctx := colly.NewContext()
		ctx.Put("filename", "gpx/routen/"+gpxfile)
		fmt.Println("got gpx download dir", gpxurl)
		c.Request("GET", e.Request.AbsoluteURL(gpxurl), nil, ctx, nil)
	})

	// route segments/POI
	c.OnHTML("#main-content > div.node.node--type-route.node--view-mode-full.ds-1col.clearfix > div.container.route-description > div > div > article", func(e *colly.HTMLElement) {
		poititle = ""
		poisubtitle = ""
		poidesc = ""
		poiimgurl = ""
		// main title
		e.ForEach("h3", func(nbr int, e *colly.HTMLElement) {
			poititle = e.Text
		})
		// subtitle
		e.ForEach("em", func(nbr int, e *colly.HTMLElement) {
			poisubtitle = strings.TrimSpace(e.Text) + "\n\n"
		})
		// POI/segment description
		e.ForEach("section > div.route-segment__description > p", func(nbr int, e *colly.HTMLElement) {
			poidesc = strings.TrimSpace(e.Text)
		})
		// image
		e.ForEach("section > div.route-segment__media > div > div", func(nbr int, e *colly.HTMLElement) {
			poiimgurl = strings.Split(e.Attr("data-flickity-bg-lazyload"), "?")[0]
		})
		contentfmt := `### %s

{{%% imgandtxt url="%s" %%}}
%s%s
{{%% /imgandtxt %%}}

`
		if poidesc != "" {
			if poiimgurl != "" {
				content += fmt.Sprintf(contentfmt, poititle, poiimgurl, poisubtitle, poidesc)
			} else {
				content += fmt.Sprintf("### %s\n\n%s%s\n\n", poititle, poisubtitle, poidesc)
			}
		}
	})

	// save gpx file
	c.OnResponse(func(r *colly.Response) {
		filename := r.Ctx.Get("filename")
		if filename != "" {
			fmt.Println("saving file", filename)
			r.Save(filename)
		}
	})

	// route page scrape -> create route .md page
	c.OnScraped(func(r *colly.Response) {
		if gpxfile != "" {
			createRoutePage()
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

func createRoutePage() {
	f, _ := os.Create("route/routen/" + name + ".md")
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
source: "be.routen"
ext_url: "%s"
gpx: "routen/%s"
length: %s
---

## Let's Go!

%s

## On Route

%s
`
	f.WriteString(fmt.Sprintf(mdContent, description, firstline(longdescription), date, description, routeurl, gpxfile, length, longdescription, content))
}

func cleanupTxt(s string) string {
	s = strings.ReplaceAll(s, "\n", "\n\n")
	s = strings.ReplaceAll(s, " ", " ")
	s = spaceStartLineRE.ReplaceAllString(s, "\n")
	s = NewLineRE.ReplaceAllString(s, "\n\n")
	s = strings.TrimSpace(s)

	return s
}

func firstline(s string) string {
	return strings.Split(strings.Split(s, ".")[0], ":")[0]
}
