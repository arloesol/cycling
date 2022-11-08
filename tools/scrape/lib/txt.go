package lib

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	spaceStartLineRE = regexp.MustCompile(`\n[ \t]*`)
	newLineRE        = regexp.MustCompile(`\n+`)
	intRE            = regexp.MustCompile(`[\d]+`)
	NodeTxtRE        = regexp.MustCompile("[ 1023456789–-]{6,}")
	NameToTitleRE    = regexp.MustCompile("[-_]")
	CamelCase        = cases.Title(language.English)
)

func CleanTxt(s string) string {
	s = strings.ReplaceAll(s, "\n", "\n\n")
	s = strings.ReplaceAll(s, " ", " ")
	s = spaceStartLineRE.ReplaceAllString(s, "\n")
	s = newLineRE.ReplaceAllString(s, "\n\n")
	s = strings.TrimSpace(s)

	return s
}

func Txtandlinks(e *colly.HTMLElement) string {
	txtfragment := CleanTxt(e.Text)
	e.ForEach("a", func(nbr int, e *colly.HTMLElement) {
		link := e.Attr("href")
		linktxt := e.Text
		if linktxt != "" {
			txtfragment = strings.Replace(txtfragment, linktxt, "["+linktxt+"]("+e.Request.AbsoluteURL(link)+")", 1)
		}
	})
	return txtfragment
}

func Firstline(s string) string {
	return strings.Split(strings.Split(s, ".")[0], ":")[0]
}

func LastSplit(s string, sep string) string {
	slice := strings.Split(s, sep)
	return slice[len(slice)-1]
}

func URLend(url string) string {
	return LastSplit(url, "/")
}

func FileExt(s string) string {
	base := strings.Split(s, "?")[0] // remove possible stuff after url
	return LastSplit(base, ".")      // ext = what's after the last "."
}

func TxttoInt(s string) int {
	intstr := intRE.FindString(s)
	if intstr != "" {
		i, err := strconv.Atoi(intstr)
		if err == nil {
			return i
		} else {
			fmt.Println("atoi error for length", err)
		}
	}
	return -1
}

func TxttoNodes(s string, route *Route) {
	for _, nodetxt := range NodeTxtRE.FindAllString(s, 5) {
		nodetxt = strings.ReplaceAll(nodetxt, "–", "-")
		nodetxt = strings.ReplaceAll(nodetxt, " ", "")
		nodes := strings.Split(nodetxt, "-")
		if len(nodes) > 3 {
			nodestr := strings.Join(nodes, ",")
			route.Nodes = append(route.Nodes, nodestr)
		}
	}
}

func NametoTitle(name string) string {
	return CamelCase.String(NameToTitleRE.ReplaceAllString(name, " "))
}
