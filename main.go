package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const sumInUabaseURL = "https://sum.in.ua"
const (
	LinkTypeIndex   = "index"
	LinkTypeArticle = "article"
)

func main() {
	parseIndex(sumInUabaseURL + "/vkazivnyk")
}

func parseIndex(url string) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	sel := doc.Find("#vkazivnyk ul li a")
	for i := range sel.Nodes {
		li := sel.Eq(i)
		href, _ := li.Attr("href")

		if isLinkAnIndex(href) {
			fmt.Printf("Index %d: %s %s\n", i, li.Text(), href)
			parseIndex(sumInUabaseURL + href)
			continue
		}

		fmt.Printf("Article %d: %s %s\n", i, li.Text(), href)
	}
}

func isLinkAnIndex(url string) bool {
	return determineLinkType(url) == LinkTypeIndex
}

func determineLinkType(url string) string {
	if strings.Contains(url, "/vkazivnyk") {
		return LinkTypeIndex
	}

	return LinkTypeArticle
}

// func main() {
// 	// Request the HTML page.
// 	res, err := http.Get("https://sum.in.ua/f/zazvychaj")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer res.Body.Close()
// 	if res.StatusCode != 200 {
// 		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
// 	}

// 	// Load the HTML document
// 	doc, err := goquery.NewDocumentFromReader(res.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Find the review items
// 	text := doc.Find("[itemprop=articleBody]").Text()
// 	fmt.Printf("Text %s\n", text)
// }
