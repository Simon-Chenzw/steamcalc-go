package main

import (
	"net/http"
	"os"
	"time"

	"steamcalc/clash"
	"steamcalc/data"
	"steamcalc/steam"

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

type MockPool struct {
	clash.ProxyPool
}

func (p *MockPool) Fetch(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

func main() {
	// data.GetDataBase()
	// cmd.Execute()

	var pool clash.ProxyPool = &MockPool{}

	item_info, err := steam.Fetch(
		pool,
		data.Item{
			Appid:    730,
			HashName: "AK-47 | Redline (Field-Tested)",
		},
	)
	if err != nil {
		panic(err)
	}

	db := data.GetDataBase()
	db.UpsertItemInfo(&item_info)
}
