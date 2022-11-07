package lib

import (
	"time"
)

type Route struct {
	Date        string // date of update
	Name        string // unique name/id
	Title       string
	Subtitle    string
	Description string   // Description of route page
	Gpxfile     string   // gpx filename
	Nodes       []string // list of nodes
	Startpunt   string   // text describing startpoint
	Content     string   // general page content
	Length      string   // length in km
	Routeurl    string   // full url of original page
	Sideimages  int      // nbr of side images
}

var Emptyroute = Route{
	Date: time.Now().UTC().Format("2006-01-02"),
}
