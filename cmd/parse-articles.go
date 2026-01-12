package cmd

import (
	"fmt"
	"sum11crawler/models"
	"sum11crawler/service"

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

		p := service.NewSumParser()
		s := service.NewSumService(db, p)

		messagesCh := make(chan string, 1)
		percentageCh := make(chan models.ParsingPercentage, 1)

		go s.ParseAndStoreArticles(messagesCh, percentageCh)

		func() {
			for {
				select {
				case mes, ok := <-messagesCh:
					if !ok {
						return
					}
					fmt.Printf("\n%s\n", mes)
				case p, ok := <-percentageCh:
					if !ok {
						return
					}
					fmt.Printf(
						"\r-------------------------\nLink: %s\nWord: %s\n%.2f%%",
						p.Word,
						p.URL,
						p.Percentage,
					)
				}
			}
		}()
	},
}
