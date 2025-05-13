package cmd

import (
	"mado/cmd/internal"
	"mado/cmd/utils"
	"mado/renderer/preview"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getGlamourStyle(theme string) glamour.TermRendererOption {
	if theme == "" {
		return glamour.WithAutoStyle()
	}

	return glamour.WithStylePath(theme)
}

var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview Markdown",
	Long:  "Renders Markdown document using glamour.",
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag(internal.ThemeVar, cmd.Flags().Lookup(internal.ThemeVar))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := utils.GetContents(viper.GetString(internal.InputFileVar), viper.GetString(internal.ReplaceFileVar))
		if err != nil {
			return err
		}

		return preview.RenderPreview(string(content), getGlamourStyle(viper.GetString(internal.ThemeVar)))
	},
}

func init() {
	previewCmd.Flags().StringVarP(&internal.Theme, internal.ThemeVar, "t", viper.GetString(internal.ThemeVar), "glamour theme or style file")

	rootCmd.AddCommand(previewCmd)
}
