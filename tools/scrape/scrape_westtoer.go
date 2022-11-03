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
	date             = ""
	name             = ""
	description      = ""
	longdescription  = ""
	length           = ""
	gpxfile          = ""
	url              = ""
	title            = ""
	startpunt        = ""
	content          = ""
	spaceStartLineRE = regexp.MustCompile("\\n[ \\t]*")
	NewLineRE        = regexp.MustCompile("\\n+")
	caser            = cases.Title(language.English)
)

func main() {
	c := colly.NewCollector()

	date = time.Now().UTC().Format("2006-01-02")

	os.Mkdir("gpx", 0750)
	os.Mkdir("gpx/westtoer", 0750)
	os.Mkdir("img", 0750)
	os.Mkdir("img/gallery", 0750)
	os.Mkdir("img/page", 0750)
	os.Mkdir("route", 0750)
	os.Mkdir("route/westtoer", 0750)

	c.OnHTML("a.node--route", func(e *colly.HTMLElement) {
		depth := e.Request.Depth
		if depth == 1 {
			src := e.Attr("href")
			fmt.Println("routeurl: ", src)
			url = e.Request.AbsoluteURL(src)
			slice := strings.Split(url, "/")
			name = "be.westtoer." + slice[len(slice)-1]
			gpxfile = ""
			content = ""
			length = ""
			startpunt = ""
			description = ""
			longdescription = ""
			title = ""
			os.Mkdir("img/gallery/"+name, 0750)
			os.Mkdir("img/page/"+name, 0750)

			e.Request.Visit(url)
		}
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		if strings.HasPrefix(e.Attr("rel"), "lightbox[field_images]") {
			src := strings.Split(e.Attr("href"), "?")[0]
			imgfilename := imgurl2imgname(src)
			if imgfilename != "" {
				ctx := colly.NewContext()
				ctx.Put("filename", "img/gallery/"+name+"/"+imgfilename)
				c.Request("GET", e.Request.AbsoluteURL(src), nil, ctx, nil)
			}
		}
	})

	c.OnHTML("meta[property=\"og:title\"]", func(e *colly.HTMLElement) {
		title = e.Attr("content")
	})

	c.OnHTML("div.field--name-field-route-distance", func(e *colly.HTMLElement) {
		e.ForEach("div.even", func(nbr int, e *colly.HTMLElement) {
			lenkm := e.Text
			length = lenkm[:len(lenkm)-2]
			length = strings.Split(length, ",")[0]
		})
	})

	c.OnHTML("div[itemprop=\"description\"]", func(e *colly.HTMLElement) {
		longdescription = cleanupTxt(e.Text)
		e.ForEach("a", func(nbr int, e *colly.HTMLElement) {
			link := e.Attr("href")
			linktxt := e.Text
			if linktxt != "" {
				longdescription = strings.Replace(longdescription, linktxt, "["+linktxt+"]("+e.Request.AbsoluteURL(link)+")", 1)
			}
		})
	})

	c.OnHTML("div.field--name-field-starting-point", func(e *colly.HTMLElement) {
		e.ForEach("div.even", func(nbr int, e *colly.HTMLElement) {
			startpunt = cleanupTxt(e.Text)
		})
	})

	c.OnHTML("div.pane-node-field-on-your-route", func(e *colly.HTMLElement) {
		content += "## En Route\n\n"
		e.ForEach("div.content", func(nbr int, e *colly.HTMLElement) { // a POI on route
			imgname := ""
			poitext := ""
			e.ForEach("h3", func(nbr int, e *colly.HTMLElement) { // title of on route POI
				content += "### " + caser.String(e.Text) + "\n\n"
			})
			e.ForEach("img[typeof=\"foaf:Image\"]", func(nbr int, e *colly.HTMLElement) {
				src := e.Attr("src")
				imgfilename := imgurl2imgname(src)
				ctx := colly.NewContext()
				ctx.Put("filename", "img/page/"+name+"/"+imgfilename)
				imgname = "/routes/page/" + name + "/" + imgfilename
				c.Request("GET", e.Request.AbsoluteURL(src), nil, ctx, nil)
			})
			e.ForEach("div.field__item > p", func(nbr int, e *colly.HTMLElement) {
				poitextfragment := e.Text
				e.ForEach("a", func(nbr int, e *colly.HTMLElement) {
					link := e.Attr("href")
					linktxt := e.Text
					if linktxt != "" {
						poitextfragment = strings.Replace(poitextfragment, linktxt, "["+linktxt+"]("+e.Request.AbsoluteURL(link)+")", 1)
					}
				})
				poitext += poitextfragment + "\n\n"
			})

			if imgname != "" {
				content += "{{% imgandtxt url=\"" + imgname + "\" %}}\n"
				content += poitext
				content += "{{% /imgandtxt %}}\n"
			} else {
				content += poitext
			}
		})
	})

	c.OnHTML("div.field--name-field-gpx", func(e *colly.HTMLElement) {
		e.ForEach("div.field__item", func(nbr int, e *colly.HTMLElement) {
			slice := strings.Split(e.Attr("resource"), "/")
			gpxfile = slice[len(slice)-1]
			gpxfile = strings.Replace(gpxfile, " ", "_", -1)
			gpxurl := ""

			e.ForEach("a", func(nbr int, e *colly.HTMLElement) {
				gpxurl = e.Attr("href")
			})

			ctx := colly.NewContext()
			ctx.Put("filename", "gpx/westtoer/"+gpxfile)
			fmt.Println("fetching gpx ", gpxfile, gpxurl)
			c.Request("GET", e.Request.AbsoluteURL(gpxurl), nil, ctx, nil)
		})
	})

	c.OnResponse(func(r *colly.Response) {
		filename := r.Ctx.Get("filename")
		fmt.Println("saving ", filename)

		if filename != "" {
			r.Save(filename)
		}
	})

	// OnScraped is called after all OnHTMLs for a webpage have been processed - if level = 2 -> create a md file with our collected info
	c.OnScraped(func(r *colly.Response) {
		if r.Request.Depth == 2 {
			fmt.Println("page scraped ", length, title, url)
			if gpxfile == "" {
				fmt.Println("no gpxfile - no route page - cleaning up")
				os.RemoveAll("img/gallery/" + name)
				os.RemoveAll("img/page/" + name)
			} else {
				f, _ := os.Create("route/westtoer/" + name + ".md")
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
source: "be.westtoer"
ext_url: "%s"
gpx: "westtoer/%s"
length: %s
---

## Let's Go !

%s

## Start 

%s 

%s
`
				description = strings.Split(longdescription, ".")[0]
				startpunt = strings.Replace(startpunt, "\n", "\n\n", -1)
				startpunt = strings.TrimRight(startpunt, " \t\n")
				content = cleanupTxt(content)
				f.WriteString(fmt.Sprintf(mdContent, title, description, date, description, url, gpxfile, length, longdescription, startpunt, content))
			}
		}
	})

	url := "https://www.westtoer.be/nl/doen/fietsroutes?field_geofield_latlon_op=10&field_geofield_latlon=&sort_by=title2&items_per_page=900&map=false&nolocation=TRUE"
	c.Visit(url)
}

func imgurl2imgname(url string) string {
	slice := strings.Split(url, "/")
	filename := slice[len(slice)-1]
	fmt.Println(url)
	filename = strings.Replace(filename, " ", "_", -1)
	baseandext := strings.Split(filename, ".")
	if len(baseandext) == 1 {
		return ""
	}
	//filename = strings.TrimRight(baseandext[0], "0123456789 _") + "." + baseandext[1]
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
