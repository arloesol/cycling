package lib

import (
	"fmt"
	"os"
	"strings"
)

func Routepage(cfg Cfg, route Route) {
	if route.Gpxfile == "" { // no gpx file - no route page on our site
		Rmdirs(cfg, route)
	}

	f, _ := os.Create("route/" + cfg.Source + "/" + route.Name + ".md")
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
	if route.Nodes != nil {
		nodestr = "\nnodetype: \"vlaams\"\nnodes: \"" + strings.Join(route.Nodes, ",") + "\""
	}
	f.WriteString(fmt.Sprintf(mdContent, route.Title, route.Description, route.Date, route.Description, route.Routeurl, route.Gpxfile, route.Length, nodestr, route.Content, route.Startpunt))
}
