package cmd

import (
	"fmt"
	"io"
	"mado/cmd/utils"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview Markdown",
	Long:  "Renders Markdown document using glamour.",
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := utils.GetContents(inputFile, replaceFile)
		if err != nil {
			return err
		}

		var output io.Writer
		if outputFile == "" {
			stdout := utils.GetStdout()
			defer stdout.Flush()
			output = stdout
		} else {
			f, err := utils.GetWriter(outputFile, force)
			if err != nil {
				return err
			}

			defer f.Close()
			output = f
		}

		out, err := glamour.RenderWithEnvironmentConfig(string(content))
		if err != nil {
			return err
		}

		fmt.Fprint(output, out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(previewCmd)
}
