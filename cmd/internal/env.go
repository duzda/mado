package internal

import "github.com/spf13/viper"

func SetEnvironment() {
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	viper.SetEnvPrefix("mado")
	viper.AutomaticEnv()

	viper.SetDefault(InputFileVar, "")
	viper.SetDefault(OutputFileVar, "")
	viper.SetDefault(ForceVar, false)
	viper.SetDefault(ReplaceFileVar, "")
	viper.SetDefault(LanguageVar, "javascript")
	viper.SetDefault(Config, "")
}

const (
	InputFileVar   = "input"
	OutputFileVar  = "output"
	ForceVar       = "force"
	ReplaceFileVar = "replace"
	LanguageVar    = "language"
	ConfigVar      = "config"
)

var (
	InputFile   string
	OutputFile  string
	Force       bool
	ReplaceFile string
	Language    string
	Config      string
)
