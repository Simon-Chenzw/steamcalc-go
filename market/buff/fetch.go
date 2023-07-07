package buff

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

const (
	BUFF_API = "https://buff.163.com/api/market/goods"
	AGENT    = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
)

func Iter(params map[string]string) <-chan Item {
	ch := make(chan Item)
	p := make(map[string]string)
	for k, v := range params {
		p[k] = v
	}

	go func() {
		slog.Info("Iterating buff...", slog.Any("params", params))

		defer close(ch)

		for page := 1; ; page++ {
			p["page_num"] = strconv.Itoa(page)
			resp, err := Fetch(p)
			if err != nil {
				slog.Error("buff.Fetch %v", err)
			}
			for _, item := range resp.Data.Items {
				ch <- item
			}
			if page >= resp.Data.TotalPage {
				break
			}
		}
	}()

	return ch
}

func Fetch(params map[string]string) (Response, error) {
	slog.Debug("buff fetch", slog.Any("params", params))

	url, err := url.Parse(BUFF_API)
	if err != nil {
		return Response{}, err
	}

	query := url.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	query.Add("game", "csgo")
	query.Add("sort_by", "price.asc")
	query.Add("use_suggestion", "0")
	query.Add("trigger", "undefined_trigger")
	query.Add("_", strconv.FormatInt(time.Now().UTC().UnixMilli(), 10))
	url.RawQuery = query.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return Response{}, err
	}

	req.Header.Add("User-Agent", AGENT)
	req.Header.Add("Cookie", viper.GetString("market.buff.cookie"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Response{}, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	var v Response
	if err := json.Unmarshal(data, &v); err != nil {
		return Response{}, err
	}

	slog.Debug("buff fetched", slog.Any("params", params), slog.Int("total_page", v.Data.TotalPage))

	return v, nil
}
