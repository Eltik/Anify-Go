package mappings

import (
	database_fetch "anify/eltik/go/src/database/impl/fetch"
	events "anify/eltik/go/src/lib"
	"anify/eltik/go/src/lib/impl/helper"
	providers "anify/eltik/go/src/mappings"
	"anify/eltik/go/src/types"
	"fmt"
	"log"
	"strings"
)

func LoadMappings(data struct {
	ID      string
	Type    types.Type
	Formats []types.Format
}) ([]types.Anime, []types.Manga, error) {
	existing, err := database_fetch.Get(data.ID, data.Type)
	if err != nil {
		log.Println("Failed to fetch existing data:", err)
		return nil, nil, err
	}

	if existing != nil {
		if existingData, ok := existing.(types.Anime); ok && existingData.ID != "" {
			events.Bus.Publish(events.COMPLETED_MAPPING_LOAD)
			return nil, nil, nil
		}
	}

	log.Println("No existing data found, fetching mappings.")

	baseProviders := providers.GetBaseProviders()

	var baseData *types.MediaInfo

	for _, provider := range *baseProviders {
		for _, format := range provider.GetFormats() {
			println("Checking format:", format)
			if format == data.Formats[0] {
				media, err := provider.GetMedia(data.ID)
				if err != nil {
					fmt.Println("Error fetching media:", err)
					return nil, nil, err
				}
				baseData = &media
				break
			}
		}
		if baseData != nil {
			break
		}
	}

	if baseData == nil || ((baseData.Title.English == nil || len(*baseData.Title.English) == 0) && (baseData.Title.Romaji == nil || len(*baseData.Title.Romaji) == 0) && (baseData.Title.Native == nil || len(*baseData.Title.Native) == 0)) {
		println("Media not found. Skipping...")

		events.Bus.Publish(events.COMPLETED_MAPPING_LOAD)
		return nil, nil, nil
	}

	var suitableProviders types.MappingsProviders

	if data.Type == types.TypeAnime {
		animeProviders := providers.GetAnimeProviders()
		for _, provider := range *animeProviders {
			for _, format := range provider.GetFormats() {
				if format == data.Formats[0] {
					suitableProviders.AnimeProviders = append(suitableProviders.AnimeProviders, provider)
					break
				}
			}
		}
	} else {
		mangaProviders := providers.GetMangaProviders()
		for _, provider := range *mangaProviders {
			for _, format := range provider.GetFormats() {
				if format == data.Formats[0] {
					suitableProviders.MangaProviders = append(suitableProviders.MangaProviders, provider)
					break
				}
			}
		}
	}

	println("Searching for media...")
	results := searchMedia(*baseData, suitableProviders)
	println("Found", len(results), "results.")

	var mappings []types.MappedResult = make([]types.MappedResult, 0)

	for _, result := range results {
		var title string
		if baseData.Title.English != nil && len(*baseData.Title.English) > 0 {
			title = *baseData.Title.English
		} else if baseData.Title.Romaji != nil && len(*baseData.Title.Romaji) > 0 {
			title = *baseData.Title.Romaji
		} else if baseData.Title.Native != nil && len(*baseData.Title.Native) > 0 {
			title = *baseData.Title.Native
		}

		providerTitles := make([][]string, 0)
		for _, r := range result {
			titles := append([]string{r.Title}, r.AltTitles...)
			filteredTitles := []string{}
			for _, title := range titles {
				if helper.IsString(title) {
					filteredTitles = append(filteredTitles, title)
				}
			}

			providerTitles = append(providerTitles, filteredTitles)
		}

		if len(providerTitles) == 0 {
			println("No titles found for " + title + " in provider " + result[0].ProviderId)
			continue
		}

		println("Found titles for provider " + result[0].ProviderId)

		titles := []string{
			*baseData.Title.English,
			*baseData.Title.Romaji,
			*baseData.Title.Native,
		}

		titles = append(titles, baseData.Synonyms...)

		var filteredTitles []string
		for _, title := range titles {
			if helper.IsString(title) {
				filteredTitles = append(filteredTitles, clean(title))
			}
		}

		bestMatchIndex := FindBestMatch2DArray(filteredTitles, providerTitles)

		if bestMatchIndex.BestMatch.Rating < 0.7 {
			//console.log(colors.gray("Unable to match ") + colors.blue(title) + colors.gray(" for ") + colors.blue(suitableProviders[i].id) + colors.gray(".") + colors.gray(" Best match rating: ") + colors.blue(bestMatchIndex.bestMatch.rating + "") + colors.gray(". ID: ") + colors.blue(providerData[bestMatchIndex.bestMatchIndex].id) + colors.gray(". Title: ") + colors.blue(providerData[bestMatchIndex.bestMatchIndex].title) + colors.gray("."));
			continue
		}

		best := result[bestMatchIndex.BestMatchIndex]

		if best.Format != types.FormatUnknown && baseData.Format != types.FormatUnknown && best.Format != baseData.Format {
			continue
		}

		if best.Year != 0 && *baseData.Year != 0 && best.Year != *baseData.Year {
			continue
		}

		var altTitles []string

		// Append non-nil Title fields to altTitles
		if baseData.Title.Romaji != nil {
			altTitles = append(altTitles, *baseData.Title.Romaji)
		}
		if baseData.Title.English != nil {
			altTitles = append(altTitles, *baseData.Title.English)
		}
		if baseData.Title.Native != nil {
			altTitles = append(altTitles, *baseData.Title.Native)
		}

		// Append synonyms
		altTitles = append(altTitles, baseData.Synonyms...)

		sim := Similarity(title, best.Title, altTitles)

		if sim.Value < 0.4 {
			continue
		}

		// if (mappings.filter((m) => m.data.id === best.id).length > 0) continue;
		if len(mappings) > 0 {
			for _, m := range mappings {
				if m.Data.ID == best.ID {
					continue
				}
			}
		}

		mappings = append(mappings, types.MappedResult{
			ID:         best.ID,
			Slug:       Slugify(best.Title),
			Data:       best,
			Similarity: sim.Value,
		})
	}

	if len(mappings) == 0 {
		println("No mappings found.")
	}

	println("Found", len(mappings), "mappings.")

	return nil, nil, nil
}

func searchMedia(baseData types.MediaInfo, suitableProviders types.MappingsProviders) [][]types.Result {
	titlesToSearch := []string{
		*baseData.Title.English,
		*baseData.Title.Romaji,
		*baseData.Title.Native,
	}

	titlesToSearch = append(titlesToSearch, baseData.Synonyms...)

	var allResults [][]types.Result

	for _, title := range titlesToSearch {
		if title == "" {
			continue
		}

		for _, provider := range suitableProviders.AnimeProviders {
			results, err := provider.Search(title, baseData.Format, *baseData.Year)
			if err != nil {
				log.Println("Error searching for anime:", err)
				continue
			}
			allResults = append(allResults, results)
		}

		for _, provider := range suitableProviders.MangaProviders {
			results, err := provider.Search(title, baseData.Format, *baseData.Year)
			if err != nil {
				log.Println("Error searching for manga:", err)
				continue
			}
			allResults = append(allResults, results)
		}
	}

	return allResults
}

func createMedia(mappings []types.MappedResult, type_ types.Type) {
	results := make([]types.Media, 0)
	for _, mapping := range mappings {
		hasPushed := false

		/*
			animeProviders := providers.GetAnimeProviders()
			mangaProviders := providers.GetMangaProviders()

			var providerType types.ProviderType

			if type_ == types.TypeAnime {
				for _, provider := range *animeProviders {
					if provider.GetID() == mapping.Data.ProviderId {
						providerType = provider.GetType()
						break
					}
				}
			} else {
				for _, provider := range *mangaProviders {
					if provider.GetID() == mapping.Data.ProviderId {
						providerType = provider.GetType()
						break
					}
				}
			}
		*/

		for _, result := range results {
			if result.ID == mapping.ID {
				hasPushed = true

				toPush := types.Mapping{
					ID:         mapping.Data.ID,
					ProviderID: mapping.Data.ProviderId,
					Similarity: mapping.Similarity,
				}

				result.Mappings = append(result.Mappings, toPush)
			}
		}

		if !hasPushed {
			data := types.Media{
				ID:   mapping.ID,
				Slug: mapping.Slug,
				Type: type_,
				Title: types.Title{
					English: nil,
					Romaji:  nil,
					Native:  nil,
				},
				Mappings: []types.Mapping{
					{
						ID:         mapping.Data.ID,
						ProviderID: mapping.Data.ProviderId,
						Similarity: mapping.Similarity,
					},
				},
				Synonyms:          []string{},
				CountryOfOrigin:   nil,
				CoverImage:        nil,
				BannerImage:       nil,
				Trailer:           nil,
				Status:            nil,
				Season:            types.SeasonUnknown,
				CurrentEpisode:    nil,
				Description:       nil,
				Duration:          nil,
				Color:             nil,
				Year:              nil,
				Rating:            nil,
				Popularity:        nil,
				AverageRating:     nil,
				AveragePopularity: nil,
				Genres:            []string{},
				Format:            types.FormatUnknown,
				Relations:         []types.Relations{},
				TotalEpisodes:     nil,
				Episodes:          types.EpisodeCollection{},
				Tags:              []string{},
				Artwork:           nil,
				Characters:        nil,
				CurrentChapter:    nil,
				TotalVolumes:      nil,
				Publisher:         nil,
				Author:            nil,
				TotalChapters:     nil,
				Chapters:          types.ChapterCollection{},
			}

			results = append(results, data)
		}
	}

	for _, media := range results {

	}
}

func clean(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
