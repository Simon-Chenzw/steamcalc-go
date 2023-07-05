package clash

import (
	"io"
	"net/http"

	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/constant"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
)

const agent = "ClashForAndroid/2.5.12.premium"

func GetProxies(url string) ([]constant.Proxy, error) {
	slog.Info("Get proxies from remote", slog.String("url", url))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", agent)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return ParseProxies(body)
}

func ParseProxies(config []byte) ([]constant.Proxy, error) {
	var data struct {
		Proxies []map[string]any `yaml:"proxies"`
	}

	err := yaml.Unmarshal(config, &data)
	if err != nil {
		return nil, err
	}

	proxies := make([]constant.Proxy, 0)

	for _, p := range data.Proxies {
		proxy, err := adapter.ParseProxy(p)
		if err != nil {
			return nil, err
		}
		proxies = append(proxies, proxy)
	}

	slog.Info("Get proxies", slog.Int("count", len(proxies)))

	return proxies, nil
}
