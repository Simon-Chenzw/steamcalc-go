package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"steamcalc/clash"

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
	proxies, err := clash.GetProxies(viper.GetString("proxy.config"))
	if err != nil {
		panic(err)
	}

	pool, err := clash.NewProxyPool(proxies, "https://ipinfo.io", 5*time.Second)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for range [10]int{} {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, "GET", "https://ipinfo.io", nil)
			if err != nil {
				slog.Error("failed to create request", slog.Any("err", err))
				return
			}

			resp, err := pool.Fetch(req)
			if err != nil {
				slog.Error("failed to fetch", slog.Any("err", err))
				return
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				slog.Error("failed to read body", slog.Any("err", err))
				return
			}
			defer resp.Body.Close()

			var v map[string]interface{}
			err = json.Unmarshal(data, &v)
			if err != nil {
				slog.Error("failed to unmarshal", slog.Any("err", err))
				return
			}
			fmt.Println(v["ip"])
		}()
	}

	wg.Wait()
	pool.Close()
}
