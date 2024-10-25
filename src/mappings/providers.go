package providers

import (
	mangaProviders "anify/eltik/go/src/mappings/impl/manga"
	manga "anify/eltik/go/src/types"
)

func GetMangaProviders() *[]manga.MangaProvider {
	providers := []manga.MangaProvider{
		mangaProviders.NewMangaDexProvider(),
	}

	return &providers
}
