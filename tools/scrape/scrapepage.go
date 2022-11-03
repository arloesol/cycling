package main

import (
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		fmt.Printf("Image found: %q -> %s\n", e.Text, link)
		c.Request("GET",e.Request.AbsoluteURL(link), nil, nil, nil)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String());
	})

	c.OnResponse(func(r *colly.Response) {
	  r.Save(r.FileName());
	})


  url := os.Args[1]
	c.Visit(url)
}
