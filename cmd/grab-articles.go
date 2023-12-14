package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sum11crawler/models"

	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func init() {
	rootCmd.AddCommand(grabArticlesCmd)
}

var grabArticlesCmd = &cobra.Command{
	Use:   "grab-articles",
	Short: "Grab articles",
	Long:  `Grab articles`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		// Migrate the schema
		db.AutoMigrate(&models.Link{})

		grabArticles(db)
	},
}

func grabArticles(db *gorm.DB) {
	var links []models.Link
	q := db.Where("html IS NULL OR html == ?", "").Order("RANDOM()")
	err := q.Find(&links).Error
	if err != nil {
		panic(err)
	}

	fmt.Printf("Selected links: %d\n", len(links))

	for i, link := range links {
		fmt.Printf("Link %d: %s\n", i, link.URL)

		// Request the HTML page.
		res, err := http.Get(sumInUaBaseURL + link.URL)
		if err != nil {
			log.Fatal(err)
		}
		// defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Printf("status code error: %d %s", res.StatusCode, res.Status)
			continue
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}

		htmlContent := string(b)
		link.HTML = &htmlContent
		db.Save(link)

		res.Body.Close()
	}
}
