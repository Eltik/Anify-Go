package types

import (
	"anify/eltik/go/src/lib/impl/request"
	"net/http"
)

type MediaInfoKeys []string

type InformationProvider[T Media, U MediaInfo] interface {
	Info(media T) (U, error)
	Request(config http.Request, proxyRequest *bool) (request.Response, error)
	GetSharedArea() MediaInfoKeys
	GetPriorityArea() MediaInfoKeys
	ProxyCheck() (bool, error)
	GetID() string
	GetType() ProviderType
}

type BaseInformationProvider struct {
	Id                 string
	Url                string
	ProviderType       ProviderType
	CustomProxy        *string
	NeedsProxy         bool
	UseGoogleTranslate bool
	OverrideProxy      bool
}

func (b *BaseInformationProvider) Info(media Media) (MediaInfo, error) {
	return MediaInfo{}, nil
}

func (b *BaseInformationProvider) GetPriorityArea() MediaInfoKeys {
	return MediaInfoKeys{}
}

func (b *BaseInformationProvider) GetSharedArea() MediaInfoKeys {
	return MediaInfoKeys{}
}

func (b *BaseInformationProvider) Request(config http.Request, proxyRequest *bool) (request.Response, error) {
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

func (b *BaseInformationProvider) ProxyCheck() (bool, error) {
	return false, nil
}

func (b *BaseInformationProvider) GetID() string {
	return b.Id
}

func (b *BaseInformationProvider) GetType() ProviderType {
	return b.ProviderType
}
