package cmd // import "iris.arke.works/forum/cmd"

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

func initConf() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("arke")
	viper.AddConfigPath("/etc/arke.d/")
	viper.AddConfigPath("$XDG_CONFIG_PATH/arke/")
	viper.AddConfigPath("$HOME/.config/arke/")
	viper.AddConfigPath(".")

	initDBConf()
	initLogConf()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error on config file: %s\n", err))
	}
}

func initDBConf() {
	viper.SetDefault("db.postgres.host", "localhost")
	viper.SetDefault("db.postgres.user", "postgres")
	viper.SetDefault("db.postgres.pass", "")
	viper.SetDefault("db.postgres.dbname", "arke")
	viper.SetDefault("db.postgres.sslmode", "verify-full")
}

func initLogConf() {
	viper.SetDefault("log.level", "warn")
	viper.BindPFlag("log.level", RootCmd.PersistentFlags().Lookup("log.level"))
}
