package service

import (
	"errors"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type SumParser struct {
}

func NewSumParser() *SumParser {
	return &SumParser{}
}

func (s *SumParser) ParseArticle(html string) (word, title, desc string, err error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return
	}

	articleSel := doc.Find("[itemprop=articleBody]")
	if articleSel == nil {
		err = errors.New("no article element")
	}

	wordSel := articleSel.Find("[itemprop=headline]")
	if wordSel != nil {
		word = wordSel.Text() // Example: АБОНУВА́ТИСЯ
	}

	// Replace acute accents
	title = strings.ReplaceAll(word, "\u0301", "") // Example: АБОНУВАТИСЯ

	desc = s.parseDesc(articleSel)
	return
}

// parseDesc gets the combined text contents of each element in the set of matched
// elements, including their descendants without new lines in Data.
func (s *SumParser) parseDesc(sel *goquery.Selection) string {
	var builder strings.Builder

	// Slightly optimized vs calling Each: no single selection object created
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			// Remove new lines
			data := strings.ReplaceAll(n.Data, "\n", " ")
			builder.WriteString(data)
		}
		// Add new lines before <p> and <br> tags
		if n.Type == html.ElementNode && (slices.Contains([]string{"p", "br"}, n.Data)) {
			builder.WriteString("\n\n")
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	for _, n := range sel.Nodes {
		f(n)
	}

	return strings.TrimSpace(builder.String())
}
