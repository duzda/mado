package cmd

import (
	"io"
	"mado/cmd/internal"
	"mado/cmd/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "Replace Markdown",
	Long:  "Replaces defined occurrences in Markdown.",
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag(internal.OutputFileVar, cmd.Flags().Lookup(internal.OutputFileVar))
		_ = viper.BindPFlag(internal.ForceVar, cmd.Flags().Lookup(internal.ForceVar))
		_ = viper.BindPFlag(internal.ReplaceFileVar, cmd.Flags().Lookup(internal.ReplaceFileVar))
	},
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

		_, err = output.Write(content)
		if err != nil {
			return err
		}

		return err
	},
}

func init() {
	replaceCmd.Flags().StringVarP(&internal.OutputFile, internal.OutputFileVar, "o", viper.GetString(internal.OutputFileVar), "file to write contents to, omitting means stdout")
	replaceCmd.Flags().BoolVarP(&internal.Force, internal.ForceVar, "f", viper.GetBool(internal.ForceVar), "overwrite existing file")
	replaceCmd.Flags().StringVarP(&internal.ReplaceFile, internal.ReplaceFileVar, "r", viper.GetString(internal.ReplaceFileVar), "path to replacements file")

	rootCmd.AddCommand(replaceCmd)
}
