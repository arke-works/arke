package cmd // import "iris.arke.works/forum/cmd"

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uber-go/zap"
)

var log zap.Logger

var RootCmd = &cobra.Command{
	Use:   "arke",
	Short: "TBD",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		lvl := zap.DynamicLevel()
		// Set Logging Level
		if viper.GetBool("verbose") {
			lvl.SetLevel(zap.DebugLevel)
		} else if viper.IsSet("log.level") {
			switch viper.GetString("log.level") {
			case "debug":
				lvl.SetLevel(zap.DebugLevel)
			case "warn":
				lvl.SetLevel(zap.WarnLevel)
			case "dpanic":
				lvl.SetLevel(zap.DPanicLevel)
			case "error":
				lvl.SetLevel(zap.ErrorLevel)
			case "fatal":
				lvl.SetLevel(zap.FatalLevel)
			case "info":
				lvl.SetLevel(zap.InfoLevel)
			case "panic":
				lvl.SetLevel(zap.PanicLevel)
			default:
				return errors.New("Unknown Logging Level")
			}
		}

		log = zap.New(zap.NewTextEncoder(), lvl)
		return nil
	},
}

func init() {
	// We must set the persistent flags *before* we start up the configuration init function
	RootCmd.PersistentFlags().String("log.level", "warn", "Set logging level")
	RootCmd.PersistentFlags().Bool("verbose", false, "Enable Verbose Mode. This overrides the logging level setting")

	initConf()
}
