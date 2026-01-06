package cmd

import (
	"sum11crawler/models"

	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// clearParsingCmd represents the clear-parsing command
var clearParsingCmd = &cobra.Command{
	Use:   "clear-parsing",
	Short: "Clear parsed data",
	Long:  `Clear parsed data`,
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
		Updates(map[string]interface{}{"word": nil, "title": nil, "desc": nil}).Error

	if err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(clearParsingCmd)
}
