package cmd

import (
	"fmt"
	"mado/cmd/internal"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "mado",
	Short: "Convert Markdown to HTML and Jira",
	Long:  "Mado is a tool to convert Markdown to HTML, Jira and preview the contents.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		configFile := viper.GetString(internal.ConfigVar)
		if configFile == "" {
			return nil
		}

		viper.SetConfigFile(configFile)
		viper.SetConfigType("env")
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Can't read config file: "+configFile)
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to mado! Use --help for usage.")
	},
}

func init() {
	internal.SetEnvironment()

	rootCmd.PersistentFlags().StringVarP(&internal.Config, internal.ConfigVar, "c", viper.GetString(internal.ConfigVar), "path to env file to be loaded")
	_ = viper.BindPFlag(internal.ConfigVar, rootCmd.PersistentFlags().Lookup(internal.ConfigVar))

	rootCmd.PersistentFlags().StringVarP(&internal.InputFile, internal.InputFileVar, "i", viper.GetString(internal.InputFileVar), "file to be processed, omit for interactive input")
	_ = viper.BindPFlag(internal.InputFileVar, rootCmd.PersistentFlags().Lookup(internal.InputFileVar))
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
