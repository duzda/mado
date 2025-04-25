package cmd

import (
	"io"

	"mado/cmd/utils"
	"mado/parser"

	"github.com/spf13/cobra"
)

var htmlCmd = &cobra.Command{
	Use:   "html",
	Short: "Convert Markdown to HTML",
	Long:  "Converts Markdown document to HTML.",
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

		parser.ToHtml(content, output)
		return err
	},
}

func init() {
	htmlCmd.Flags().StringVarP(&replaceFile, "replace", "r", "", "file with replaces to be used")

	rootCmd.AddCommand(htmlCmd)
}
