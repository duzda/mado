package cmd

import (
	"io"
	"mado/cmd/utils"

	"github.com/spf13/cobra"
)

var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "Replace Markdown",
	Long:  "Replaces defined occurrences in Markdown.",
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

		_, err = output.Write(content)
		if err != nil {
			return err
		}

		return err
	},
}

func init() {
	replaceCmd.Flags().StringVarP(&replaceFile, "replace", "r", "", "file with replaces to be used")

	rootCmd.AddCommand(replaceCmd)
}
