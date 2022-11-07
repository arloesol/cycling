package lib

import (
	"time"
)

type RoutePOI struct {
	Title     string
	Imgurl    string      // The image related to the POI - only 1 image - should be external url - no download
	Content   string      // Content in pure txt
	ContLinks [][2]string // Anchorlinks related to Content : 0:txt and 1:url
	Extlink   string      // External weblink with more info about POI
}
type Route struct {
	Date        string      // date of update
	Shortname   string      // unique id for the route source
	Name        string      // full unique id : source id + Shortname
	Title       string      // Human readable title
	Subtitle    string      // Other short info
	Description string      // Description of route page
	Gpxfile     string      // gpx filename
	Nodes       []string    // list of nodes - nodestring is like "32,44,71,33"
	Startpunt   string      // text describing startpoint
	Content     string      // general page content
	ContLinks   [][2]string // Anchorlinks related to Content : 0:txt and 1:url
	Length      int         // length in km
	Routeurl    string      // full url of original page
	Sideimages  int         // nbr of side images
	Tags        []string    // extra tags
	Categories  []string    // extra categories
	POIs        []RoutePOI  // Points Of Interest on route
	Signage     string      // Signage of the route
}

var Emptyroute = Route{
	Date: time.Now().UTC().Format("2006-01-02"),
}
