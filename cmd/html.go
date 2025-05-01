package cmd

import (
	"io"

	"mado/cmd/internal"
	"mado/cmd/utils"
	"mado/parser"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var htmlCmd = &cobra.Command{
	Use:   "html",
	Short: "Convert Markdown to HTML",
	Long:  "Converts Markdown document to HTML.",
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := utils.GetContents(viper.GetString(internal.InputFileVar), viper.GetString(internal.ReplaceFileVar))
		if err != nil {
			return err
		}

		var output io.Writer
		outputFile := viper.GetString(internal.OutputFileVar)
		if outputFile == "" {
			stdout := utils.GetStdout()
			defer utils.JoinErrors(&err, stdout.Flush)
			output = stdout
		} else {
			f, err := utils.GetWriter(outputFile, viper.GetBool(internal.ForceVar))
			if err != nil {
				return err
			}

			defer utils.JoinErrors(&err, f.Close)
			output = f
		}

		parser.ToHtml(content, output)
		return err
	},
}

func init() {
	htmlCmd.Flags().StringVarP(&internal.ReplaceFile, internal.ReplaceFileVar, "r", viper.GetString(internal.ReplaceFileVar), "file to write contents to, omitting means stdout")
	_ = viper.BindPFlag(internal.ReplaceFileVar, htmlCmd.Flags().Lookup(internal.ReplaceFileVar))

	rootCmd.AddCommand(htmlCmd)
}
