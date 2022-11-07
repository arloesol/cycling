package lib

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

var (
	spaceStartLineRE = regexp.MustCompile(`\n[ \t]*`)
	NewLineRE        = regexp.MustCompile(`\n+`)
)

func CleanTxt(s string) string {
	s = strings.ReplaceAll(s, "\n", "\n\n")
	s = strings.ReplaceAll(s, " ", " ")
	s = spaceStartLineRE.ReplaceAllString(s, "\n")
	s = NewLineRE.ReplaceAllString(s, "\n\n")
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
