package request

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	proxies "example/user/hello/src/lib/impl/proxies"
)

func GetRandomUnbannedProxy(providerId string) *string {
	var data []proxies.Proxy

	switch {
	case providerId == "novelupdates":
		data = proxies.MangaProxies
	case contains(providerId, []string{"base1"}):
		data = proxies.BaseProxies
	case contains(providerId, []string{"anime1"}):
		data = proxies.AnimeProxies
	case contains(providerId, []string{"manga1", "mangadex"}):
		data = proxies.MangaProxies
	default:
		data = proxies.MetaProxies
	}

	if len(data) == 0 {
		return nil
	}

	randomProxy := data[rand.Intn(len(data))]
	return &randomProxy.IP
}

// contains checks if a provider ID exists in a list
func contains(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func Request(providerId string, useGoogleTranslate bool, config http.Request, proxyRequest bool) (*http.Response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	copyHeaders := func(src, dst http.Header) {
		for key, values := range src {
			for _, value := range values {
				dst.Add(key, value)
			}
		}
	}

	if proxyRequest {
		if useGoogleTranslate {
			encodedURL := url.QueryEscape(config.URL.String())

			// Modify the request to use Google Translate.
			translatedURL := fmt.Sprintf("http://translate.google.com/translate?sl=ja&tl=en&u=%s", encodedURL)
			req, err := http.NewRequestWithContext(context.Background(), config.Method, translatedURL, config.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %w", err)
			}

			// Copy headers from the original request.
			req.Header = config.Header

			// Send the request via Google Translate.
			resp, err := client.Do(req)
			if err != nil {
				return nil, fmt.Errorf("request through Google Translate failed: %w", err)
			}
			return resp, nil
		}

		proxy := GetRandomUnbannedProxy(providerId)
		if proxy == nil {
			return nil, fmt.Errorf("no unbanned proxy available for provider: %s", providerId)
		}

		proxiedURL := fmt.Sprintf("%s/%s", *proxy, config.URL.String())
		req, err := http.NewRequestWithContext(context.Background(), config.Method, proxiedURL, config.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Origin", config.URL.String())
		copyHeaders(config.Header, req.Header)

		// Send the request via the unbanned proxy.
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request through proxy failed: %w", err)
		}

		return resp, nil
	} else {
		req, err := http.NewRequestWithContext(context.Background(), config.Method, config.URL.String(), config.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Send the request directly.
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		return resp, nil
	}
}

// Response example
type Response struct {
	Ok       bool   `json:"ok"`
	Status   int    `json:"status"`
	Text     string `json:"text"`
	Response *http.Response
}
