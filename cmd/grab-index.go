package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sum11crawler/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

const (
	indexMain   = "/vkazivnyk"      // СУМ-11
	indexNewest = "/vkazivnyk/vits" // Новітній користувацькій словник
	nameMain    = "sum.in.ua"
	nameNewest  = "sum.in.ua_user"
)

func init() {
	rootCmd.AddCommand(grabIndexCmd)

	grabIndexCmd.Flags().StringVarP(&vocabulary, "vocabulary", "v", "sum.in.ua", "Vocabulary")
}

var vocabulary string

var grabIndexCmd = &cobra.Command{
	Use:   "grab-index",
	Short: "Grab index",
	Long:  `Grab index`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		// Migrate the schema
		db.AutoMigrate(&models.Link{})

		// Index URL
		index := indexMain
		if vocabulary != nameMain {
			index = indexNewest
		}

		// Name
		name := nameMain
		if vocabulary != nameMain {
			name = nameNewest
		}

		parseIndex(db, sumInUaBaseURL+index, name, nil)
	},
}

func parseIndex(db *gorm.DB, url string, vocabulary string, parentLink *models.Link) {
	if db == nil {
		panic("parse error. no DB")
	}

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
			// Create index
			fmt.Printf("Index %d: %s %s\n", i, li.Text(), href)

			newIndex := models.Link{
				URL:        href,
				Type:       LinkTypeIndex,
				Vocabulary: vocabulary,
			}
			if parentLink != nil {
				newIndex.ParentID = &parentLink.ID
				newIndex.ParentURL = &parentLink.URL
			}
			db.Create(&newIndex)

			parseIndex(db, sumInUaBaseURL+href, vocabulary, &newIndex)
			continue
		}

		// Create article
		fmt.Printf("Article %d: %s %s\n", i, li.Text(), href)

		newArticle := models.Link{
			URL:        href,
			Type:       LinkTypeArticle,
			Vocabulary: vocabulary,
		}
		if parentLink != nil {
			newArticle.ParentID = &parentLink.ID
			newArticle.ParentURL = &parentLink.URL
		}
		db.Create(&newArticle)
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
