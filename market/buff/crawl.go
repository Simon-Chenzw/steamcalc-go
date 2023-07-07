package buff

import (
	"strconv"

	"steamcalc/data"

	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

func Crawl(min_price, max_price float64) {
	ch := make(chan Item, 1000)

	min_price_str := strconv.FormatFloat(min_price, 'f', 2, 64)
	max_price_str := strconv.FormatFloat(max_price, 'f', 2, 64)

	go func() {
		defer close(ch)

		for _, category := range viper.GetStringSlice("market.buff.category") {
			for item := range Iter(
				map[string]string{
					"category":  category,
					"min_price": min_price_str,
					"max_price": max_price_str,
				},
			) {
				ch <- item
			}
		}

		for _, group := range viper.GetStringSlice("market.buff.category_group") {
			for item := range Iter(
				map[string]string{
					"category_group": group,
					"min_price":      min_price_str,
					"max_price":      max_price_str,
				},
			) {
				ch <- item
			}
		}
	}()

	db := data.GetDataBase()
	for item := range ch {
		slog.Debug("buff crawled", slog.String("name", item.MarketHashName), slog.Float64("price", item.SellMinPrice))

		p := data.Item{
			Appid:    uint32(item.Appid),
			HashName: item.MarketHashName,
			MarketPrice: []data.MarketPrice{{
				Source: "buff",
				Price:  item.SellMinPrice,
			}},
		}

		db.UpsertItem(&p)
	}
}
