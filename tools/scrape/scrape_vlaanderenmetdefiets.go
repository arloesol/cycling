package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var (
  date = "";
  name = "";
  description = "";
  length = "";
  gpxfile = "";
  url = "";
)

func main() {
	c := colly.NewCollector();
	
	date = time.Now().UTC().Format("2006-01-02");
	
	os.Mkdir("gpx", 0750);
	os.Mkdir("gpx/vlaanderen", 0750);
	os.Mkdir("img", 0750);
	os.Mkdir("route", 0750);
	os.Mkdir("route/vlaanderen", 0750);

	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		depth := e.Request.Depth;
		
		if (depth == 1) {
  		if (strings.HasPrefix(e.Text, "Bekijk ")) {
  		  //fmt.Printf("link found: %q -> %s\n", e.Text, link);
  		  slice := strings.Split(link, "/");
  		  name = "be.vlaanderenmetdefiets." + slice[len(slice)-1];
  		  os.Mkdir("img/" + name, 0750);
  		  url = link;
  		  e.Request.Visit(link); // check the route page
  		};
  	};

		if (depth == 2) {
		  if (strings.Contains(e.Text, "GPX")) {
		    filename := e.Attr("download");
		    filename = strings.Replace(filename, ".GPX", ".gpx", -1);
		    filename = strings.Replace(filename, " ", "_", -1);
		    fmt.Println("GPX: ", filename);
		    gpxfile = filename;
  		  ctx := colly.NewContext();
  		  ctx.Put("filename", "gpx/vlaanderen/" + filename);
		    c.Request("GET",e.Request.AbsoluteURL(link), nil, ctx, nil)
		  }
  	};
	})
	
	c.OnHTML("meta", func(e *colly.HTMLElement) {
		metaname := e.Attr("name")
		if (metaname == "description") {
		  description = e.Attr("content");
		  fmt.Println("route: ", description);
		};
	});
		
	c.OnHTML("span", func(e *colly.HTMLElement) {
		depth := e.Request.Depth;
		if (depth == 2) {
  		class := e.Attr("class")
  		if (class == "field-distance__value") {
  		  length = e.Text;
  		  fmt.Println("length: ", length);
  		};
		};
	});

	c.OnHTML("img", func(e *colly.HTMLElement) {
		depth := e.Request.Depth;
		if (depth == 2) {
  		src := e.Attr("src")
  		if (strings.Contains(src, "images")) {
  		  filename := imgurl2imgname(src, e.Attr("alt"));
  		  ctx := colly.NewContext();
  		  ctx.Put("filename", "img/" + name + "/" + filename);
  		  fmt.Println("image: ", src, filename, e.Attr("alt"));
		    c.Request("GET",e.Request.AbsoluteURL(src), nil, ctx, nil)
  		};
		};
	});

  c.OnResponse(func(r *colly.Response) {
    filename := r.Ctx.Get("filename");
    if (filename != "") {
	    r.Save(filename);
    }
	})
	
	// OnScraped is called after all OnHTMLs for a webpage have been processed - if level = 2 -> create a md file with our collected info
	c.OnScraped(func(r *colly.Response) {
	  if (r.Request.Depth == 2) {
	    fmt.Println("page scraped ", description)
	    f, _ := os.Create("route/vlaanderen/" + name + ".md");
	    defer f.Close()
	    mdContent := `---
title: "!! TITLE !!"
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
source: "be.vlaanderenmetdefiets"
ext_url: "%s"
gpx: "vlaanderen/%s"
length: %s
---
`;	    
      lengthnokm := length[:len(length)-2]
	    f.WriteString(fmt.Sprintf(mdContent, description, date, description, url, gpxfile, lengthnokm));
	  }
	})

  url := "https://www.vlaanderenmetdefiets.be/"
	c.Visit(url)
}

func imgurl2imgname (url string, alt string) string {
	filename := strings.Replace(alt," ", "_", -1);
	filename = strings.Split(filename, "_©")[0];
  slice := strings.Split(url, ".");
  ext := slice[len(slice)-1];
  return filename + "." + ext;
}