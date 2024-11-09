package types

import (
	"anify/eltik/go/src/lib/impl/request"
	"net/http"
)

type SeasonalResponse struct {
	Seasonal []MediaInfo `json:"seasonal"`
	Trending []MediaInfo `json:"trending"`
	Popular  []MediaInfo `json:"popular"`
	Top      []MediaInfo `json:"top"`
}

type ScheduleResponse struct {
	Sunday    []MediaInfo `json:"sunday"`
	Monday    []MediaInfo `json:"monday"`
	Tuesday   []MediaInfo `json:"tuesday"`
	Wednesday []MediaInfo `json:"wednesday"`
	Thursday  []MediaInfo `json:"thursday"`
	Friday    []MediaInfo `json:"friday"`
	Saturday  []MediaInfo `json:"saturday"`
}

type BaseProvider interface {
	Search(query string, mediaType Type, formats []Format, page int, perPage int) ([]MediaInfo, error)
	SearchAdvanced(query string, mediaType Type, formats []Format, page int, perPage int, genres []string, genresExcluded []string, season Season, year int, tags []string, tagsExcluded []string) ([]MediaInfo, error)
	GetCurrentSeason() (Season, error)
	GetMedia(id string) (MediaInfo, error)
	GetSeasonal(mediaType Type, formats []Format) (SeasonalResponse, error)
	GetSchedule() (ScheduleResponse, error)
	GetIds() ([]string, error)
	Request(config http.Request, proxyRequest *bool) (request.Response, error)
	ProxyCheck() (bool, error)
	GetFormats() []Format
}

type BaseBaseProvider struct {
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

func (b *BaseBaseProvider) Search(query string, mediaType Type, formats []Format, page int, perPage int) ([]MediaInfo, error) {
	return nil, nil
}

func (b *BaseBaseProvider) SearchAdvanced(query string, mediaType Type, formats []Format, page int, perPage int, genres []string, genresExcluded []string, season Season, year int, tags []string, tagsExcluded []string) ([]MediaInfo, error) {
	return nil, nil
}

func (b *BaseBaseProvider) GetCurrentSeason() (Season, error) {
	return SeasonUnknown, nil
}

func (b *BaseBaseProvider) GetMedia(id string) (MediaInfo, error) {
	return MediaInfo{}, nil
}

func (b *BaseBaseProvider) GetSeasonal(mediaType Type, formats []Format) (SeasonalResponse, error) {
	return SeasonalResponse{}, nil
}

func (b *BaseBaseProvider) GetSchedule() (ScheduleResponse, error) {
	return ScheduleResponse{}, nil
}

func (b *BaseBaseProvider) GetIds() ([]string, error) {
	return nil, nil
}

func (b *BaseBaseProvider) GetFormats() []Format {
	return b.Formats
}

func (b *BaseBaseProvider) Request(config http.Request, proxyRequest *bool) (request.Response, error) {
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

func (b *BaseBaseProvider) ProxyCheck() (bool, error) {
	return false, nil
}
