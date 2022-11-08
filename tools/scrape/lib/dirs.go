package lib

import "os"

var (
	datadir       string
	sourcedir     string
	gpxdir        string
	imgdir        string
	imggallerydir string
	routedir      string
	logdir        string
)

func Mkalldirs(cfg Cfg) {
	os.Mkdir(datadir, 0750)
	sourcedir = datadir + "/" + cfg.Source
	gpxdir = sourcedir + "/gpx"
	imgdir = sourcedir + "/img"
	imggallerydir = imgdir + "/gallery"
	routedir = sourcedir + "/route"
	os.Mkdir(sourcedir, 0750)
	os.Mkdir(gpxdir, 0750)
	os.Mkdir(imgdir, 0750)
	os.Mkdir(imggallerydir, 0750)
	os.Mkdir(routedir, 0750)
}

func Mkdirs(cfg Cfg, route Route) {
	os.Mkdir(imggallerydir+"/"+route.Name, 0750)
}

func Rmdirs(cfg Cfg, route Route) {
	os.RemoveAll(imggallerydir + "/" + route.Name)
}
