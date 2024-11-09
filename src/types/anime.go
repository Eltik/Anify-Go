package types

import (
	"anify/eltik/go/src/lib/impl/request"
	"net/http"
)

type AnimeProvider interface {
	Search(query string, format Format, year int) ([]Result, error)
	FetchEpisodes(id string) ([]Episode, error)
	FetchRecent() ([]Anime, error)
	FetchSources(id string, proxy bool, chapter *Chapter) (interface{}, error) // can return []Page or string
	Request(config http.Request, proxyRequest *bool) (request.Response, error)
	ProxyCheck() (bool, error)
	GetFormats() []Format
	GetID() string
	GetType() ProviderType
}

type BaseAnimeProvider struct {
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

func (b *BaseAnimeProvider) Search(query string, format Format, year int) ([]Result, error) {
	return nil, nil
}

func (b *BaseAnimeProvider) FetchEpisodes(id string) ([]Chapter, error) {
	return nil, nil
}

func (b *BaseAnimeProvider) FetchRecent() ([]Manga, error) {
	return nil, nil
}

func (b *BaseAnimeProvider) FetchSources(id string, proxy bool, chapter *Chapter) (interface{}, error) {
	return nil, nil
}

func (b *BaseAnimeProvider) Request(config http.Request, proxyRequest *bool) (request.Response, error) {
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

func (b *BaseAnimeProvider) ProxyCheck() (bool, error) {
	return false, nil
}

func (b *BaseAnimeProvider) GetFormats() []Format {
	return b.Formats
}

func (b *BaseAnimeProvider) GetID() string {
	return b.Id
}

func (b *BaseAnimeProvider) GetType() ProviderType {
	return b.ProviderType
}
