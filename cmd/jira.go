package cmd

import (
	"io"
	"mado/parser"
	"os"

	"github.com/spf13/cobra"
)

var (
	language string
)

var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "Convert Markdown to Jira",
	Long:  "Converts Markdown document to Jira specific format.",
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := os.ReadFile(inputFile)
		if err != nil {
			return err
		}

		var output io.Writer
		if outputFile == "" {
			stdout := getStdout()
			defer stdout.Flush()
			output = stdout
		} else {
			f, err := getWriter(outputFile, force)
			if err != nil {
				return err
			}

			defer f.Close()
			output = f
		}

		parser.ToJira(content, output, language)
		return nil
	},
}

func init() {
	jiraCmd.Flags().StringVarP(&language, "language", "l", "javascript", "programming language to be used for code blocks")

	rootCmd.AddCommand(jiraCmd)
}
