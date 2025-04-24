package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	inputFile   string
	outputFile  string
	force       bool
	replaceFile string
)

var rootCmd = &cobra.Command{
	Use:   "mado",
	Short: "Convert Markdown to HTML and Jira",
	Long:  "Mado is a tool to convert Markdown to HTML, Jira and preview the contents.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to mado! Use --help for usage.")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "", "file to be processed")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "file to write contents to, omitting means stdout")
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "overwrite existing file")
	rootCmd.MarkFlagRequired("input")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
