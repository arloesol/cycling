package lib

import (
	"fmt"
	"os"
	"strings"
)

func extraYaml(cfg Cfg, route Route) string {
	return tags(cfg, route) + categories(cfg, route) + nodes(cfg, route)
}

func tags(cfg Cfg, route Route) string { // tags : cfg and route tags + tag depending on length
	// Todo: add a tag based on slope - flat, rolling, hilly, ... - check official categories
	alltags := append(cfg.Tags, route.Tags...)
	return "\ntags:\n - \"" + strings.Join(alltags, "\"\n - \"") + "\""
}

func categories(cfg Cfg, route Route) string { // categories : cfg and route categories + "route"
	allcategories := append(append(cfg.Categories, "route"), route.Categories...)
	return "\ncategories:\n - \"" + strings.Join(allcategories, "\"\n - \"") + "\""
}

func nodes(cfg Cfg, route Route) string {
	if route.Nodes != nil {
		return "\nnodetype: \"" + cfg.NodeType + "\"\nnodes:\n - \"" + strings.Join(route.Nodes, "\"\n - \"") + "\""
	} else {
		return ""
	}
}

func poicontent(route Route) string {
	content := ""
	for _, poi := range route.POIs {
		// title
		content += "### " + CamelCase.String(poi.Title) + "\n\n"
		// include links in poitext in markdown format
		poimd := poi.Content
		for _, link := range poi.ContLinks {
			poimd = strings.Replace(poimd, link[0], "["+link[0]+"]("+link[1]+")", 1)
		}
		// use hugo imgandtxt shortcode in case of poi image
		if poi.Imgurl != "" {
			content += "{{% imgandtxt url=\"" + poi.Imgurl + "\""
			if poi.Extlink != "" {
				content += " extlink=\"" + poi.Extlink + "\""
			}
			if poi.ImgAlt != "" {
				content += " alt=\"" + cleanAttr(poi.ImgAlt) + "\""
			} else {
				content += " alt=\"" + cleanAttr(poi.Title) + "\""
			}
			content += " %}}\n"
			content += strings.TrimSpace(poimd) + "\n"
			content += "{{% /imgandtxt %}}\n\n"
		} else { // Todo: use the external link even if no poi image
			content += poimd
		}
	}
	if content != "" {
		content = "\n\n## On route\n\n" + content
	}
	return content
}

func Routepage(cfg Cfg, route Route) {
	if route.Gpxfile == "" { // no gpx file - no route page on our site
		LogWarning.Println("no gpx - skipping route", route.Name)
		Rmdirs(cfg, route)
	} else {

		f, _ := os.Create(routedir + "/" + route.Name + ".md")
		defer f.Close()
		mdContent := `---
date: "%s"
title: "%s"
subtitle: "%s"
description: "%s"
region: "%s"
website: "%s"
ext_url: "%s"%s
routes:
    - name: Main
      gpx: "%s"
      length: %d000
      up: 0
      down: 0
      minheight: 99999999999999
      maxheight: 0
      minslope: 0
      maxslope: 0
      avgposslope: 0
      avgnegslope: 0
      slopehisto:
        - 0
        - 0
        - 0
        - 0
        - 0
      effortlevel: 0
      minlat: 0
      minlon: 0
      maxlat: 0
      maxlon: 0
---

## Let's Go ! 

%s
`
		if route.Subtitle == "" {
			route.Subtitle = Firstline(route.Description)
		}
		route.Content = strings.TrimSpace(route.Content)
		if route.Startpunt != "" {
			route.Content += "\n\n## Start\n\n" + route.Startpunt
		}
		if route.Signage != "" {
			route.Content += "\n\n## Signage\n\n" + route.Signage
		}
		if route.POIs != nil {
			route.Content += poicontent(route)
		}
		LogInfo.Println("Create route page", route.Name)
		f.WriteString(fmt.Sprintf(mdContent,
			route.Date, route.Title, route.Subtitle, route.Description,
			cfg.Region, cfg.Srcpfx[:len(cfg.Srcpfx)-1],
			route.Routeurl, extraYaml(cfg, route),
			cfg.Source+"/"+route.Gpxfile, route.Length, route.Content))
	}
}

func cleanAttr(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\"", "")
}
