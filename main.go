package main

import (
	"os"
	"time"

	"steamcalc/cmd"
	"steamcalc/data"

	"github.com/lmittmann/tint"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

func init() { // slog
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))
}

func init() { // viper
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	data.GetDataBase()
	cmd.Execute()
}
