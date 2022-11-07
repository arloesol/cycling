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

func SaveGPX(c *colly.Collector, e *colly.HTMLElement, cfg Cfg, route *Route) {
	route.Gpxfile = route.Shortname + ".gpx"
	if cfg.Savegpx {
		ctx := colly.NewContext()
		ctx.Put("filename", "gpx/"+cfg.Source+"/"+route.Gpxfile)
		c.Request("GET", e.Request.AbsoluteURL(e.Attr("href")), nil, ctx, nil)
	}
}

// image saving - image filename - filename parameter
//   "0" -> imgname = routename
//   "1" -> imgname = routename_n
//   Attr-attr -> based on e.Attr("attr") value
//   else just use filename value

func SaveIMGanchor(c *colly.Collector, e *colly.HTMLElement, cfg Cfg, route *Route, attr string, filename string) {
	imgurl := e.Request.AbsoluteURL(e.Attr(attr))
	if strings.HasPrefix(filename, "Attr-") {
		filename = imgattr2imgname(imgurl, e.Attr(filename[5:]))
		fmt.Println("anchor", e.Attr(filename[5:]), "-", filename)
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
			ctx.Put("filename", "img/gallery/"+route.Name+"/"+imgfilename)
			c.Request("GET", imgurl, nil, ctx, nil)
		}
	}
}

func save(cfg Cfg, r *colly.Response, filename string) {
	fmt.Println("saving", filename)
	r.Save(filename)
}

func SaveOnResponse(cfg Cfg) func(*colly.Response) {
	return func(r *colly.Response) {
		filename := r.Ctx.Get("filename")
		if filename != "" {
			save(cfg, r, filename)
		}
		filename2 := r.Ctx.Get("filename2")
		if filename2 != "" {
			save(cfg, r, filename2)
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
		filename = route.Title + "_" + fmt.Sprintf("%d", route.Sideimages) + "." + ext
		route.Sideimages += 1
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
	base = strings.Split(base, "_©")[0] // flandersbybike stuff
	return base + "." + ext
}
