package cmd

import (
	"io"
	"mado/cmd/utils"
	"mado/parser"

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
		content, err := utils.GetContents(inputFile, replaceFile)
		if err != nil {
			return err
		}

		var output io.Writer
		if outputFile == "" {
			stdout := utils.GetStdout()
			defer func() {
				derr := stdout.Flush()
				if err == nil {
					err = derr
				}
			}()
			output = stdout
		} else {
			f, err := utils.GetWriter(outputFile, force)
			if err != nil {
				return err
			}

			defer func() {
				derr := f.Close()
				if err == nil {
					err = derr
				}
			}()
			output = f
		}

		parser.ToJira(content, output, language)
		return nil
	},
}

func init() {
	jiraCmd.Flags().StringVarP(&language, "language", "l", "javascript", "programming language to be used for code blocks")
	jiraCmd.Flags().StringVarP(&replaceFile, "replace", "r", "", "file with replaces to be used")

	rootCmd.AddCommand(jiraCmd)
}
