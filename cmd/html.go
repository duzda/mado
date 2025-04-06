package cmd

import (
	"io"
	"os"

	"mado/parser"

	"github.com/spf13/cobra"
)

var htmlCmd = &cobra.Command{
	Use:   "html",
	Short: "Convert Markdown to HTML",
	Long:  "Converts Markdown document to HTML.",
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := os.ReadFile(input_file)
		if err != nil {
			return err
		}

		var output io.Writer
		if output_file == "" {
			stdout := getStdout()
			defer stdout.Flush()
			output = stdout
		} else {
			f, err := getWriter(output_file, force)
			if err != nil {
				return err
			}

			defer f.Close()
			output = f
		}

		parser.ToHtml(content, output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(htmlCmd)
}
