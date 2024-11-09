package types

import (
	"anify/eltik/go/src/lib/impl/request"
	"net/http"
	"strings"
)

type MangaProvider interface {
	Search(query string, format Format, year int) ([]Result, error)
	FetchChapters(id string) ([]Chapter, error)
	FetchRecent() ([]Manga, error)
	FetchPages(id string, proxy bool, chapter *Chapter) (interface{}, error) // can return []Page or string
	Request(config http.Request, proxyRequest *bool) (request.Response, error)
	ProxyCheck() (bool, error)
	PadNum(number string, places int) string
	GetFormats() []Format
	GetID() string
	GetType() ProviderType
}

type BaseMangaProvider struct {
	RateLimit          int
	Id                 string
	Url                string
	Formats            []Format
	ProviderType       ProviderType
	CustomProxy        *string
	NeedsProxy         bool
	UseGoogleTranslate bool
	OverrideProxy      bool
}

func (b *BaseMangaProvider) Search(query string, format Format, year int) ([]Result, error) {
	return nil, nil
}

func (b *BaseMangaProvider) FetchChapters(id string) ([]Chapter, error) {
	return nil, nil
}

func (b *BaseMangaProvider) FetchRecent() ([]Manga, error) {
	return nil, nil
}

func (b *BaseMangaProvider) FetchPages(id string, proxy bool, chapter *Chapter) (interface{}, error) {
	return nil, nil
}

func (b *BaseMangaProvider) Request(config http.Request, proxyRequest *bool) (request.Response, error) {
	if proxyRequest == nil {
		proxyRequest = &b.NeedsProxy
	}
	if *proxyRequest && !b.NeedsProxy {
		*proxyRequest = false
	}

	resp, err := request.Request(b.Id, b.UseGoogleTranslate, config, *proxyRequest)
	if err != nil {
		return request.Response{}, err
	}

	return request.Response{Response: resp}, nil
}

func (b *BaseMangaProvider) ProxyCheck() (bool, error) {
	return false, nil
}

func (b *BaseMangaProvider) PadNum(number string, places int) string {
	parts := strings.Split(number, "-")
	for i, part := range parts {
		part = strings.TrimSpace(part)
		digits := len(strings.Split(part, ".")[0])
		padding := strings.Repeat("0", max(0, places-digits))
		parts[i] = padding + part
	}
	return strings.Join(parts, "-")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (b *BaseMangaProvider) GetFormats() []Format {
	return b.Formats
}

func (b *BaseMangaProvider) GetID() string {
	return b.Id
}

func (b *BaseMangaProvider) GetType() ProviderType {
	return b.ProviderType
}
