package information

import (
	"anify/eltik/go/src/lib/impl/request"
	"anify/eltik/go/src/types"
	"net/http"
)

type MangaDexInformationProvider struct {
	types.BaseInformationProvider
	Api string
}

func NewMangaDexInformationProvider() *MangaDexInformationProvider {
	return &MangaDexInformationProvider{
		BaseInformationProvider: types.BaseInformationProvider{
			Id:                 "mangadex",
			Url:                "https://mangadex.org",
			ProviderType:       types.ProviderTypeManga,
			NeedsProxy:         true,
			UseGoogleTranslate: false,
		},
		Api: "https://api.mangadex.org",
	}
}

func (p *MangaDexInformationProvider) Info(media types.Media) (types.MediaInfo, error) {
	return types.MediaInfo{}, nil
}

func (p *MangaDexInformationProvider) GetSharedArea() types.MediaInfoKeys {
	return types.MediaInfoKeys{"synonyms", "genres", "artwork", "tags"}
}

func (p *MangaDexInformationProvider) ProxyCheck() (bool, error) {
	return false, nil
}

func (p *MangaDexInformationProvider) Request(config http.Request, proxyRequest *bool) (request.Response, error) {
	return p.BaseInformationProvider.Request(config, proxyRequest)
}
