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
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := utils.GetContents(viper.GetString(internal.InputFileVar), viper.GetString(internal.ReplaceFileVar))
		if err != nil {
			return err
		}

		return preview.RenderPreview(string(content), getGlamourStyle(viper.GetString(internal.ThemeVar)))
	},
}

func init() {
	previewCmd.PersistentFlags().StringVarP(&internal.Theme, internal.ThemeVar, "t", viper.GetString(internal.ThemeVar), "glamour theme or style file")
	_ = viper.BindPFlag(internal.ThemeVar, previewCmd.PersistentFlags().Lookup(internal.ThemeVar))

	rootCmd.AddCommand(previewCmd)
}
