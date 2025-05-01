package cmd

import (
	"io"
	"mado/cmd/internal"
	"mado/cmd/utils"
	"mado/parser"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "Convert Markdown to Jira",
	Long:  "Converts Markdown document to Jira specific format.",
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

		parser.ToJira(content, output, viper.GetString(internal.LanguageVar))
		return err
	},
}

func init() {
	jiraCmd.Flags().StringVarP(&internal.Language, internal.LanguageVar, "l", viper.GetString(internal.LanguageVar), "programming language to be used for code blocks")
	_ = viper.BindPFlag(internal.LanguageVar, jiraCmd.Flags().Lookup(internal.LanguageVar))

	jiraCmd.Flags().StringVarP(&internal.ReplaceFile, internal.ReplaceFileVar, "r", viper.GetString(internal.ReplaceFileVar), "file to write contents to, omitting means stdout")
	_ = viper.BindPFlag(internal.ReplaceFileVar, jiraCmd.Flags().Lookup(internal.ReplaceFileVar))

	rootCmd.AddCommand(jiraCmd)
}
