package clash

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/Dreamacro/clash/constant"
)

func urlToMetadata(u *url.URL) (constant.Metadata, error) {
	port := u.Port()
	if port == "" {
		switch u.Scheme {
		case "https":
			port = "443"
		case "http":
			port = "80"
		default:
			return constant.Metadata{}, fmt.Errorf("%s scheme not Support", u)
		}
	}

	return constant.Metadata{
		Host:    u.Hostname(),
		DstIP:   nil,
		DstPort: port,
	}, nil
}

func Fetch(proxy constant.Proxy, req *http.Request) (*http.Response, error) {
	meta, err := urlToMetadata(req.URL)
	if err != nil {
		return nil, err
	}

	dial, err := proxy.DialContext(req.Context(), &meta)
	if err != nil {
		return nil, err
	}
	defer dial.Close()

	transport := &http.Transport{
		Dial: func(string, string) (net.Conn, error) { return dial, nil },
		// from http.DefaultTransport
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	defer client.CloseIdleConnections()

	return client.Do(req)
}
