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

		s.ParseAndStoreArticles(
			func(mes string) {
				fmt.Printf("\n%s\n", mes)
			},
			func(p models.ParsingPercentage) {
				fmt.Printf(
					"\r-------------------------\nLink: %s\nWord: %s\n%.2f%%",
					p.Word,
					p.URL,
					p.Percentage,
				)
			},
		)
	},
}
