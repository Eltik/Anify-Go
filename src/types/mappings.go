package types

type MappingsProviders struct {
	AnimeProviders []AnimeProvider `json:"animeProviders"`
	MangaProviders []MangaProvider `json:"mangaProviders"`
}

type MappedResult struct {
	ID         string  `json:"id"`
	Slug       string  `json:"slug"`
	Data       Result  `json:"data"`
	Similarity float64 `json:"similarity"`
}

type StringRating struct {
	Target string
	Rating float64
}

type StringResult struct {
	Ratings        []StringRating
	BestMatch      StringRating
	BestMatchIndex int
}

type ToPush struct {
	ID           string       `json:"id"`
	ProviderID   string       `json:"provider_id"`
	ProviderType ProviderType `json:"provider_type"`
	Similarity   float64      `json:"similarity"`
}
