package lib

import "os"

func Mkalldirs(cfg Cfg) {
	os.Mkdir("gpx", 0750)
	os.Mkdir("gpx/"+cfg.Source, 0750)
	os.Mkdir("img", 0750)
	os.Mkdir("img/gallery", 0750)
	os.Mkdir("route", 0750)
	os.Mkdir("route/"+cfg.Source, 0750)
}

func Mkdirs(cfg Cfg, route Route) {
	os.Mkdir("img/gallery/"+route.Name, 0750)
}

func Rmdirs(cfg Cfg, route Route) {
	os.RemoveAll("img/gallery/" + route.Name)
}
