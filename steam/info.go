package steam

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"steamcalc/clash"
	"steamcalc/data"

	"golang.org/x/exp/slog"
)

const (
	AGENT = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
)

var /* const */ (
	ItemNameidPattern = regexp.MustCompile(`Market_LoadOrderSpread\( (\d*) \);`)
	LocaleNamePattern = regexp.MustCompile(`<a href="https://steamcommunity.com/market/listings/\d+/[^"]+">([^<]+)</a>`)
)

func Fetch(pool *clash.ProxyPool, item data.Item) (data.ItemInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := fmt.Sprintf(
		"https://steamcommunity.com/market/listings/%d/%s",
		item.Appid,
		url.PathEscape(item.HashName),
	)
	slog.Debug("Steam Fetch", slog.String("url", url))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return data.ItemInfo{}, err
	}
	req.Header.Set("User-Agent", AGENT)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,en;q=0.7")

	resp, err := pool.Fetch(req)
	if err != nil {
		return data.ItemInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return data.ItemInfo{}, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return data.ItemInfo{}, err
	}

	item_nameid, err := strconv.Atoi(ItemNameidPattern.FindStringSubmatch(string(html))[1])
	if err != nil {
		return data.ItemInfo{}, err
	}

	locale_name := LocaleNamePattern.FindStringSubmatch(string(html))[1]

	return data.ItemInfo{
		Item:       item,
		ItemNameid: item_nameid,
		LocaleName: locale_name,
	}, nil
}
