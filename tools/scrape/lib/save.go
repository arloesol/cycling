package lib

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func SaveGPX(c *colly.Collector, e *colly.HTMLElement, cfg Cfg, route *Route) {
	route.Gpxfile = route.Title + ".gpx"
	ctx := colly.NewContext()
	ctx.Put("filename", "gpx/"+cfg.Source+"/"+route.Gpxfile)
	c.Request("GET", e.Request.AbsoluteURL(e.Attr("href")), nil, ctx, nil)
}

func SaveIMGanchor(c *colly.Collector, e *colly.HTMLElement, cfg Cfg, route *Route, side byte) {
	// side '' ->  no side management
	// side '0' -> main img
	// side '1' -> side img
	imgurl := e.Request.AbsoluteURL(e.Attr("src"))
	SaveIMG(c, imgurl, cfg, route, side)
}

func SaveIMG(c *colly.Collector, imgurl string, cfg Cfg, route *Route, side byte) {
	imgfilename := imgname(imgurl, route, side)
	if imgfilename != "" {
		ctx := colly.NewContext()
		ctx.Put("filename", "img/gallery/"+route.Name+"/"+imgfilename)
		c.Request("GET", imgurl, nil, ctx, nil)
	}
}

func SaveOnResponse(r *colly.Response) {
	filename := r.Ctx.Get("filename")
	if filename != "" {
		fmt.Println("saving file", filename)
		r.Save(filename)
	}
}

func imgname(url string, route *Route, side byte) string {
	slice := strings.Split(url, "/")
	filename := slice[len(slice)-1]
	slice = strings.Split(filename, ".")
	if len(slice) == 1 {
		return ""
	}
	ext := slice[len(slice)-1]
	if side == '1' {
		filename = route.Title + "_" + fmt.Sprintf("%d", route.Sideimages) + "." + ext
		route.Sideimages += 1
	} else {
		filename = route.Title + "." + ext
	}

	return filename
}
