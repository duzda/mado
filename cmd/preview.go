package cmd

import (
	"fmt"
	"io"
	"mado/cmd/internal"
	"mado/cmd/utils"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview Markdown",
	Long:  "Renders Markdown document using glamour.",
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

		out, err := glamour.RenderWithEnvironmentConfig(string(content))
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(output, out)
		if err != nil {
			return err
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(previewCmd)
}
