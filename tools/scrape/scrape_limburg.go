package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	nodes            = []string{}
	startpunt        = ""
	content          = ""
	routeurl         = ""
	sideimages       = 0
	spaceStartLineRE = regexp.MustCompile("\\n[ \\t]*")
	NewLineRE        = regexp.MustCompile("\\n+")
	NodesRE          = regexp.MustCompile(" |\\.")
	caser            = cases.Title(language.English)
	poiurl           = ""
	poialt           = ""
	poiimgurl        = ""
	poititle         = ""
	poidesc          = ""
)

func main() {
	c := colly.NewCollector()

	date = time.Now().UTC().Format("2006-01-02")

	os.Mkdir("gpx", 0750)
	os.Mkdir("gpx/limburg", 0750)
	os.Mkdir("img", 0750)
	os.Mkdir("img/gallery", 0750)
	os.Mkdir("route", 0750)
	os.Mkdir("route/limburg", 0750)

	// route page urls from overview page
	c.OnHTML("div.overview-search-results", func(e *colly.HTMLElement) {
		e.ForEach("a", func(nbr int, e *colly.HTMLElement) {
			routeurl = e.Request.AbsoluteURL(e.Attr("href"))
			if strings.HasPrefix(routeurl, "https://www.visitlimburg.be/en/route") {
				fmt.Println("visiting", routeurl)
				slice := strings.Split(routeurl, "/")
				title = slice[len(slice)-1]
				name = "be.visitlimburg." + title
				routenbr += 1
				content = ""
				gpxfile = ""
				os.Mkdir("img/gallery/"+name, 0750)
				c.Request("GET", routeurl, nil, nil, nil)
			}
		})
	})

	// route description
	c.OnHTML("#block-sassy-content > article > header > div > div > div > div:nth-child(2) > div > div > div > div.slide-content-left.col-md-5.col-md-offset-2 > h1", func(e *colly.HTMLElement) {
		description = cleanupTxt(e.Text)
		description = strings.ReplaceAll(description, "\"", "")
		//fmt.Println(description)
	})
	// route long description
	c.OnHTML("#block-sassy-content > article > div > div > div.left > div.route--description", func(e *colly.HTMLElement) {
		longdescription = cleanupTxt(e.Text)
		//fmt.Println(longdescription)
	})

	// route length
	c.OnHTML("#block-sassy-content > article > div > div > div.right > div.route--info > div.route--info-card > ul > li:nth-child(1) > span:nth-child(2)", func(e *colly.HTMLElement) {
		length = e.Text[:len(e.Text)-3]
	})

	// startpunt
	c.OnHTML("#block-sassy-content > article > div > div > div.right > div.route--info > div.route--info-card > ul > li:nth-child(3) > span:nth-child(2)", func(e *colly.HTMLElement) {
		startpunt = cleanupTxt(e.Text)
		//fmt.Println("startpunt", startpunt)
	})

	// gpx file
	c.OnHTML("#block-sassy-content > article > div > div > div.right > div.route--info > div.route--info-download > ul > li:nth-child(1) > a", func(e *colly.HTMLElement) {
		gpxurl := e.Attr("href")
		gpxfile = title + ".gpx"
		ctx := colly.NewContext()
		ctx.Put("filename", "gpx/limburg/"+gpxfile)
		// fmt.Println("fetching gpx ", gpxfile, gpxurl)
		c.Request("GET", e.Request.AbsoluteURL(gpxurl), nil, ctx, nil)
	})

	// main image
	c.OnHTML("#block-sassy-content > article > header > div > div > div > div:nth-child(1) > div > div > img", func(e *colly.HTMLElement) {
		src := e.Attr("src")
		imgfile := "img/gallery/" + name + "/" + imgname(src)
		ctx := colly.NewContext()
		ctx.Put("filename", imgfile)
		//fmt.Println("main image ", imgfile, src)
		c.Request("GET", e.Request.AbsoluteURL(src), nil, ctx, nil)
	})
	// on route POIs
	c.OnHTML("#block-sassy-content > article > footer > div.field--name-route-poi > section > section > div > a", func(e *colly.HTMLElement) {
		// fmt.Println("a POI found")
		poiurl = e.Request.AbsoluteURL(e.Attr("href"))
		poialt = ""
		poiimgurl = ""
		poititle = ""
		poidesc = ""
		e.ForEach("div.accomodation-image > div > img", func(nbr int, e *colly.HTMLElement) {
			poialt = e.Attr("alt")
			poiimgurl = strings.Split(e.Attr("data-src"), "?")[0]
		})
		e.ForEach("div.accomodation-title", func(nbr int, e *colly.HTMLElement) {
			poititle = cleanupTxt(e.Text)
		})
		e.ForEach("div.accomodation-description", func(nbr int, e *colly.HTMLElement) {
			poidesc = cleanupTxt(e.Text)
		})
		// fmt.Println(poititle, poidesc, poialt, poiimgurl, poiurl)
		contentfmt := `### %s

{{%% imgandtxt url="%s" extlink="%s" %%}}
%s
{{%% /imgandtxt %%}}

`
		content += fmt.Sprintf(contentfmt, poititle, poiimgurl, poiurl, poidesc)
	})

	// save gpx or image files
	c.OnResponse(func(r *colly.Response) {
		filename := r.Ctx.Get("filename")
		if filename != "" {
			fmt.Println("saving file", filename)
			r.Save(filename)
		}
	})

	// route page scrape -> create route .md page
	c.OnScraped(func(r *colly.Response) {
		if strings.HasPrefix(r.Request.URL.String(), "https://www.visitlimburg.be/en/route/") {
			fmt.Println("routepage scraped")
			createRoutePage()
		}
	})

	// visit all route ovv pages until there are no more (routes)
	oldnbrroutes := -1
	for page := 0; routenbr > oldnbrroutes; page++ {
		oldnbrroutes = routenbr
		ovvurl := "https://www.visitlimburg.be/en/find-cycling-routes?page=" + fmt.Sprintf("%d", page)
		fmt.Println(ovvurl)
		c.Visit(ovvurl)
	}

}

func createRoutePage() {
	f, _ := os.Create("route/limburg/" + name + ".md")
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
source: "be.visitlimburg"
ext_url: "%s"
gpx: "limburg/%s"
length: %s
---

## Let's Go!

%s

## Start

%s

## On Route

%s
`
	if gpxfile == "" {
		// skip this route - cleanup
		fmt.Println("no gpxfile -> no route page -> cleaning up images")
		os.RemoveAll("img/gallery/" + name)
	} else {
		f.WriteString(fmt.Sprintf(mdContent, description, firstline(longdescription), date, description, routeurl, gpxfile, length, longdescription, startpunt, content))
	}
}

func imgname(url string) string {
	slice := strings.Split(url, "/")
	filename := slice[len(slice)-1]
	filename = strings.Split(filename, "?")[0]
	slice = strings.Split(filename, ".")
	if len(slice) == 1 {
		return ""
	}
	//base := strings.Join(slice[:len(slice)-1],".")
	ext := slice[len(slice)-1]
	filename = title + "." + ext

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

func firstline(s string) string {
	return strings.Split(s, ".")[0]
}
