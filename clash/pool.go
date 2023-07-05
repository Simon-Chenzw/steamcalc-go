package clash

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/Dreamacro/clash/constant"
	"golang.org/x/exp/slog"
)

type Request struct {
	request  *http.Request
	response chan<- struct {
		*http.Response
		error
	}
}

type ProxyPool struct {
	proxies []constant.Proxy
	ch      chan Request
}

func filter(proxies []constant.Proxy, url string, threshold time.Duration) []constant.Proxy {
	slog.Info("fitler proxies with url", slog.String("url", url), slog.Duration("threshold", threshold))

	var wg sync.WaitGroup
	ch := make(chan constant.Proxy)

	for _, proxy := range proxies {
		wg.Add(1)
		go func(proxy constant.Proxy) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), threshold)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return
			}

			resp, err := Fetch(proxy, req)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == 200 {
				ch <- proxy
			}
		}(proxy)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var result []constant.Proxy
	for proxy := range ch {
		result = append(result, proxy)
	}
	slog.Info("fitler proxies with done", slog.Int("count", len(result)))
	return result
}

func NewProxyPool(proxies []constant.Proxy, test_url string, delay_threshold time.Duration) (*ProxyPool, error) {
	slog.Info("New proxy pool")

	live_proxies := filter(proxies, test_url, delay_threshold)

	if len(live_proxies) == 0 {
		return nil, errors.New("no live proxies")
	}

	p := &ProxyPool{
		proxies: proxies,
		ch:      make(chan Request),
	}

	for _, proxy := range p.proxies {
		go func(proxy constant.Proxy) {
			for req := range p.ch {
				resp, err := Fetch(proxy, req.request)
				req.response <- struct {
					*http.Response
					error
				}{resp, err}
			}
		}(proxy)
	}

	return p, nil
}

func (p *ProxyPool) Fetch(req *http.Request) (*http.Response, error) {
	future := make(chan struct {
		*http.Response
		error
	}, 1)
	p.ch <- Request{
		request:  req,
		response: future,
	}
	result := <-future
	return result.Response, result.error
}

func (p *ProxyPool) Close() {
	close(p.ch)
}
