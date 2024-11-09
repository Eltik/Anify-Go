package providers

import (
	baseProviders "anify/eltik/go/src/mappings/impl/base"
	mangaProviders "anify/eltik/go/src/mappings/impl/manga"
	types "anify/eltik/go/src/types"
)

func GetBaseProviders() *[]types.BaseProvider {
	providers := []types.BaseProvider{
		baseProviders.NewMangaDexBaseProvider(),
	}

	return &providers
}

func GetAnimeProviders() *[]types.AnimeProvider {
	providers := []types.AnimeProvider{}

	return &providers
}

func GetMangaProviders() *[]types.MangaProvider {
	providers := []types.MangaProvider{
		mangaProviders.NewMangaDexProvider(),
	}

	return &providers
}
