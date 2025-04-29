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
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := utils.GetContents(viper.GetString(internal.InputFileVar), viper.GetString(internal.ReplaceFileVar))
		if err != nil {
			return err
		}

		var output io.Writer
		outputFile := viper.GetString(internal.OutputFileVar)
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
			f, err := utils.GetWriter(outputFile, viper.GetBool(internal.ForceVar))
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

		_, err = output.Write(content)
		if err != nil {
			return err
		}

		return err
	},
}

func init() {
	replaceCmd.Flags().StringVarP(&internal.ReplaceFile, internal.ReplaceFileVar, "r", viper.GetString(internal.ReplaceFileVar), "file to write contents to, omitting means stdout")
	_ = viper.BindPFlag(internal.ReplaceFileVar, replaceCmd.Flags().Lookup(internal.ReplaceFileVar))

	rootCmd.AddCommand(replaceCmd)
}
