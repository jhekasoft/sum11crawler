package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const sumInUaBaseURL = "https://sum.in.ua"
const (
	LinkTypeIndex   = "index"
	LinkTypeArticle = "article"
)

var dbPath string
var note string

var rootCmd = &cobra.Command{
	Use:   "sum11crawler",
	Short: "SUM-11 crawler",
	Long:  `SUM-11 crawler`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dbPath, "db", "", "db.db", "Path of DB-file")
}
