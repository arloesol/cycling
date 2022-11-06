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

const ( // set to "" true true to parse everything
	pagetoparse = "" // "" to parse all routes
	savegpx     = true
	saveimg     = true
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
	bewegwijzering   = ""
	content          = ""
	spaceStartLineRE = regexp.MustCompile("\\n[ \\t]*")
	NewLineRE        = regexp.MustCompile("\\n+")
	NodesRE          = regexp.MustCompile(" |\\.")
	NodeTxtRE        = regexp.MustCompile(": [ 1023456789–]+")
	caser            = cases.Title(language.English)
	poiimgurl        = ""
	poilink          = ""
	nodes            = []string{}
	signimgurl       = ""
	signimgname      = ""
)

func main() {
	c := colly.NewCollector()

	date = time.Now().UTC().Format("2006-01-02")

	os.Mkdir("gpx", 0750)
	os.Mkdir("gpx/westtoer", 0750)
	os.Mkdir("img", 0750)
	os.Mkdir("img/gallery", 0750)
	os.Mkdir("img/signage", 0750)
	os.Mkdir("route", 0750)
	os.Mkdir("route/westtoer", 0750)

	c.OnHTML("a.node--route", func(e *colly.HTMLElement) {
		depth := e.Request.Depth
		if depth == 1 {
			src := e.Attr("href")
			if pagetoparse == "" {
				fmt.Println("routeurl: ", src)
			}
			url = e.Request.AbsoluteURL(src)
			slice := strings.Split(url, "/")
			name = "be.westtoer." + slice[len(slice)-1]
			gpxfile = ""
			content = ""
			length = ""
			startpunt = ""
			bewegwijzering = ""
			description = ""
			longdescription = ""
			title = ""
			signimgname = ""
			os.Mkdir("img/gallery/"+name, 0750)

			if pagetoparse == "" || url == pagetoparse {
				e.Request.Visit(url)
			}
		}
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		if strings.HasPrefix(e.Attr("rel"), "lightbox[field_images]") {
			src := strings.Split(e.Attr("href"), "?")[0]
			imgfilename := imgurl2imgname(src)
			if imgfilename != "" {
				ctx := colly.NewContext()
				ctx.Put("filename", "img/gallery/"+name+"/"+imgfilename)
				if saveimg {
					c.Request("GET", e.Request.AbsoluteURL(src), nil, ctx, nil)
				}
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
		// image in description and text in description "Volg deze bordjes:" -> signage
		e.ForEach("img", func(nbr int, e *colly.HTMLElement) {
			if strings.Contains(longdescription, "Volg deze bordjes:") {
				longdescription = cleanupTxt(strings.Replace(longdescription, "Volg deze bordjes:", "", 1))
				signimgurl = e.Request.AbsoluteURL(e.Attr("src"))
				signimgname = imgurl2imgname(signimgurl)
				signimgname = strings.Replace(signimgname, "bordjes-", "", -1)
				ctx := colly.NewContext()
				ctx.Put("filename", "img/signage/westtoer."+signimgname)
				ctx.Put("filename2", "img/gallery/"+name+"/"+signimgname)
				if saveimg {
					c.Request("GET", signimgurl, nil, ctx, nil)
				}
			}
		})
	})

	// startpunt
	c.OnHTML("div.field--name-field-starting-point", func(e *colly.HTMLElement) {
		e.ForEach("div.even", func(nbr int, e *colly.HTMLElement) {
			startpunt = cleanupTxt(e.Text)
		})
	})

	// bewegwijzering
	c.OnHTML("div.field--name-field-junctions", func(e *colly.HTMLElement) {
		e.ForEach("div.even", func(nbr int, e *colly.HTMLElement) {
			bewegwijzering = cleanupTxt(e.Text)
		})
		// vb Knooppunten met kasseistrook: 75 – 1 – 69 – 57 – 28 – 12 – 62 – 29 – 65 – 61 – 31 – 92 – 32 – 97 – 35 – 36 – 43 – 23 – 86 – 63 – 39 – 84 – 94 – 16 – 99 – 47 – 13 – 10 – 68 – 61 – 9 – 18 – 59 – 1 – 75
		//    Knooppunten zonder kasseistrook: 75 – 1 – 69 – 57 – 28 – 12 – 62 – 29 – 65 – 61 – 31 – 92 – 32 – 97 – 35 – 36 – 43 – 23 – 86 – 63 – 39 – 84 – 94 – 16 – 99 – 47 – 13 – 10 – 68 – 61 – 92 – 59 – 1 – 75
		// Todo: only first list of nodes -- manage more than 1 later
		nodetxt := NodeTxtRE.FindString(bewegwijzering)
		if nodetxt != "" {
			nodetxt = nodetxt[2:]
			nodetxt = NodesRE.ReplaceAllString(nodetxt, "") // remove " " and "."
			nodes = strings.Split(nodetxt, "–")
		}
	})

	c.OnHTML("div.pane-node-field-on-your-route", func(e *colly.HTMLElement) {
		content += "## En Route\n\n"
		e.ForEach("div.content", func(nbr int, e *colly.HTMLElement) { // a POI on route
			poiimgurl = ""
			poitext := ""
			poilink = ""
			e.ForEach("h3", func(nbr int, e *colly.HTMLElement) { // title of on route POI
				content += "### " + caser.String(e.Text) + "\n\n"
			})
			e.ForEach("img[typeof=\"foaf:Image\"]", func(nbr int, e *colly.HTMLElement) {
				poiimgurl = e.Request.AbsoluteURL(e.Attr("src"))
			})
			e.ForEach("div.field__item > p", func(nbr int, e *colly.HTMLElement) {
				poitextfragment := strings.TrimSpace(e.Text)
				e.ForEach("a", func(nbr int, e *colly.HTMLElement) {
					link := e.Attr("href")
					if poilink == "" {
						poilink = link
					}
					linktxt := e.Text
					if linktxt != "" {
						poitextfragment = strings.Replace(poitextfragment, linktxt, "["+linktxt+"]("+e.Request.AbsoluteURL(link)+")", 1)
					}
				})
				poitext += poitextfragment + "\n\n"
			})

			if poiimgurl != "" {
				content += "{{% imgandtxt url=\"" + poiimgurl + "\""
				if poilink != "" {
					content += " extlink=\"" + poilink + "\""
				}
				content += " %}}\n"
				content += strings.TrimSpace(poitext) + "\n"
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
			if savegpx {
				c.Request("GET", e.Request.AbsoluteURL(gpxurl), nil, ctx, nil)
			}
		})
	})

	c.OnResponse(func(r *colly.Response) {
		filename := r.Ctx.Get("filename")
		if filename != "" {
			fmt.Println("saving ", filename)
			r.Save(filename)
			filename2 := r.Ctx.Get("filename2")
			if filename2 != "" {
				fmt.Println("saving ", filename2)
				r.Save(filename2)
			}
		}
	})

	// OnScraped is called after all OnHTMLs for a webpage have been processed - if level = 2 -> create a md file with our collected info
	c.OnScraped(func(r *colly.Response) {
		if r.Request.Depth == 2 {
			fmt.Println("page scraped ", length, title, url)
			if gpxfile == "" {
				fmt.Println("no gpxfile - no route page - cleaning up")
				os.RemoveAll("img/gallery/" + name)
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
length: %s%s%s
---

## Let's Go !

%s

## Start 

%s%s%s`
				description = strings.Split(longdescription, ".")[0]
				startpunt = strings.Replace(startpunt, "\n", "\n\n", -1)
				startpunt = strings.TrimRight(startpunt, " \t\n")
				content = cleanupTxt(content)
				if content != "" {
					content = "\n\n" + content
				}
				if signimgname != "" {
					signimgname = "\nsignage: \"" + signimgname + "\""
				}
				if bewegwijzering != "" {
					bewegwijzering = "\n\n## Signage\n\n" + bewegwijzering
				}
				nodestr := ""
				if nodes != nil {
					nodestr = "\nnodetype: \"vlaams\"\nnodes: \"" + strings.Join(nodes, ",") + "\""
				}
				f.WriteString(fmt.Sprintf(mdContent, title, description, date, description, url, gpxfile, length, nodestr, signimgname, longdescription, startpunt, bewegwijzering, content))
			}
		}
	})

	url := "https://www.westtoer.be/nl/doen/fietsroutes?field_geofield_latlon_op=10&field_geofield_latlon=&sort_by=title2&items_per_page=900&map=false&nolocation=TRUE"
	c.Visit(url)
}

func imgurl2imgname(url string) string {
	slice := strings.Split(url, "/")
	filename := slice[len(slice)-1]
	filename = strings.Replace(filename, " ", "_", -1)
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
