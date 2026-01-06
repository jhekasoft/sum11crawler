package cmd

import (
	"sum11crawler/models"

	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// cleanParsingCmd represents the clean-parsing command
var cleanParsingCmd = &cobra.Command{
	Use:   "clean-parsing",
	Short: "Clean parsed data",
	Long:  `Clean parsed data`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		// Migrate the schema
		db.AutoMigrate(&models.Link{})

		clearParsing(db)
	},
}

func clearParsing(db *gorm.DB) {
	err := db.Model(&models.Link{}).
		Where("type = ?", LinkTypeArticle).
		Where(
			db.Session(&gorm.Session{}).Unscoped().
				Where("word IS NOT NULL").
				Or("title IS NOT NULL").
				Or("desc IS NOT NULL"),
		).
		Updates(map[string]any{"word": nil, "title": nil, "desc": nil}).Error

	if err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(cleanParsingCmd)
}
