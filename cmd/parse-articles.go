package cmd

import (
	"errors"
	"fmt"
	"strings"
	"sum11crawler/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
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
	var items []models.Link
	q := db.Where("(html IS NOT NULL OR html != ?) AND desc IS NULL AND type = ?", "", LinkTypeArticle)
	err := q.Find(&items).Error
	if err != nil {
		panic(err)
	}

	fmt.Printf("Selected articles: %d\n", len(items))

	for i, item := range items {
		fmt.Printf("Link %d: %s\n", i, item.URL)

		if item.HTML == nil {
			fmt.Println("No HTML")
			continue
		}
		word, title, desc, err := parseArticle(*item.HTML)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(title)

		item.Word = &word
		item.Title = &title
		item.Desc = &desc
		item.Title = &title
		db.Save(item)
	}
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
	title = strings.Replace(word, "\u0301", "", -1) // Example: АБОНУВАТИСЯ

	desc = articleSel.Text()
	return
}
