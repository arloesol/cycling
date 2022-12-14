package lib

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func Routename(cfg Cfg, route *Route, name string) {
	route.Shortname = name
	route.Name = cfg.Srcpfx + name
}

func SaveGPX(c *colly.Collector, e *colly.HTMLElement, cfg Cfg, route *Route, url string) {
	route.Gpxfile = route.Shortname + ".gpx"
	if cfg.Savegpx {
		if strings.HasPrefix(url, "Attr-") {
			url = e.Request.AbsoluteURL(e.Attr(url[5:]))
		}
		ctx := colly.NewContext()
		ctx.Put("filename", gpxdir+"/"+route.Gpxfile)
		c.Request("GET", url, nil, ctx, nil)
	}
}

// image saving - image filename - filename parameter
//   "0" -> imgname = routename
//   "1" -> imgname = routename_n
//   Attr-attr -> based on e.Attr("attr") value
//   else just use filename value

func SaveIMGanchor(c *colly.Collector, e *colly.HTMLElement, cfg Cfg, route *Route, urlattr string, filename string) {
	imgurl := e.Request.AbsoluteURL(e.Attr(urlattr))
	if strings.HasPrefix(filename, "Attr-") {
		filename = imgattr2imgname(imgurl, e.Attr(filename[5:]))
	}
	SaveIMG(c, imgurl, cfg, route, filename)
}

func SaveIMG(c *colly.Collector, imgurl string, cfg Cfg, route *Route, imgfilename string) {
	if cfg.Saveimg {
		if len(imgfilename) == 1 {
			imgfilename = imgname(imgurl, route, imgfilename)
		}
		if imgfilename != "" {
			ctx := colly.NewContext()
			ctx.Put("filename", imggallerydir+"/"+route.Name+"/"+imgfilename)
			c.Request("GET", imgurl, nil, ctx, nil)
		}
	}
}

func SaveOnResponse(cfg Cfg) func(*colly.Response) {
	return func(r *colly.Response) {
		filename := r.Ctx.Get("filename")
		if filename != "" {
			LogInfo.Println("saving", filename)
			r.Save(filename)
		}
	}
}

func imgname(url string, route *Route, side string) string {
	slice := strings.Split(url, "/")
	filename := slice[len(slice)-1]
	slice = strings.Split(filename, ".")
	if len(slice) == 1 {
		return ""
	}
	ext := slice[len(slice)-1]
	if side == "1" {
		filename = route.Title + "_" + fmt.Sprintf("%d", route.sideimages) + "." + ext
		route.sideimages += 1
	} else {
		filename = route.Title + "." + ext
	}

	return filename
}

func imgattr2imgname(url string, attr string) string {
	// using attr for base of filename and url for extension
	ext := FileExt(url)
	base := URLend(attr)
	base = strings.Split(base, "?")[0]
	base = strings.Split(base, ".")[0]
	base = strings.Replace(base, " ", "_", -1)
	base = strings.Split(base, "_??")[0] // flandersbybike stuff
	return base + "." + ext
}
