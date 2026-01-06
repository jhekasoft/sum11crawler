package cmd

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sum11crawler/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
	"gorm.io/gorm"
)

func init() {
	rootCmd.AddCommand(parseArticlesCmd)
}

var parseArticlesCmd = &cobra.Command{
	Use:   "parse-articles",
	Short: "Parse articles",
	Long:  `Parse articles`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		// Migrate the schema
		db.AutoMigrate(&models.Link{})

		parseAndStoreArticles(db)
	},
}

func parseAndStoreArticles(db *gorm.DB) {
	var total int64
	db.Model(&models.Link{}).Where("(html IS NOT NULL OR html != ?) AND type = ?", "", LinkTypeArticle).Count(&total)

	var items []models.Link
	q := db.Where("(html IS NOT NULL OR html != ?) AND desc IS NULL AND type = ?", "", LinkTypeArticle)
	err := q.Find(&items).Error
	if err != nil {
		panic(err)
	}

	itemsCount := len(items)
	doneCount := total - int64(itemsCount)

	fmt.Printf("Selected articles: %d\n", itemsCount)

	for i, item := range items {
		if item.HTML == nil {
			fmt.Println("No HTML")
			continue
		}
		word, title, desc, err := parseArticle(*item.HTML)
		if err != nil {
			fmt.Println(err.Error())
		}

		percentage := calcPercentage(int(doneCount)+i+1, int(total))

		// Output to the CLI
		fmt.Printf("\rLink: %s\nWord: %s\n%.2f%%", item.URL, word, percentage)

		item.Word = &word
		item.Desc = &desc
		item.Title = &title
		db.Save(item)
	}
}

func calcPercentage(part, total int) float32 {
	return float32(part) / float32(total) * 100
}

func parseArticle(html string) (word, title, desc string, err error) {
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

	desc = parseDesc(articleSel)
	return
}

// parseDesc gets the combined text contents of each element in the set of matched
// elements, including their descendants without new lines in Data.
func parseDesc(s *goquery.Selection) string {
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
	for _, n := range s.Nodes {
		f(n)
	}

	return strings.TrimSpace(builder.String())
}
