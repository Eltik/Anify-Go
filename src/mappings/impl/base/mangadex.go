package base

import (
	"anify/eltik/go/src/lib/impl/helper"
	"anify/eltik/go/src/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type MangaDexBaseProvider struct {
	types.BaseBaseProvider
	Api string
}

func NewMangaDexBaseProvider() *MangaDexBaseProvider {
	return &MangaDexBaseProvider{
		BaseBaseProvider: types.BaseBaseProvider{
			RateLimit:          250,
			Id:                 "mangadex",
			Url:                "https://mangadex.org",
			Formats:            []types.Format{types.FormatManga, types.FormatOneShot},
			ProviderType:       types.ProviderTypeManga,
			NeedsProxy:         true,
			UseGoogleTranslate: false,
		},
		Api: "https://api.mangadex.org",
	}
}

func (p *MangaDexBaseProvider) Search(query string, mediaType types.Type, formats []types.Format, page int, perPage int) ([]types.MediaInfo, error) {
	var results []types.MediaInfo

	for page := 0; page <= 1; page++ {
		uri, _ := url.Parse(p.Api + "/manga")
		q := uri.Query()

		q.Set("title", query)
		q.Set("limit", "25")
		q.Set("offset", strconv.Itoa(25*page))
		q.Set("order[relevance]", "desc")
		q.Add("contentRating[]", "safe")
		q.Add("contentRating[]", "suggestive")
		q.Add("includes[]", "cover_art")
		uri.RawQuery = q.Encode()

		resp, err := p.Request(http.Request{
			URL:    uri,
			Method: "GET",
		}, &p.NeedsProxy)
		if err != nil {
			return nil, err
		}
		defer resp.Response.Body.Close()

		if resp.Response.StatusCode != 200 {
			return nil, fmt.Errorf("unexpected status code: %d", resp.Response.StatusCode)
		}

		if resp.Response.Header.Get("Content-Type") != "application/json" {
			return nil, fmt.Errorf("invalid content type: %s", resp.Response.Header.Get("Content-Type"))
		}

		body, err := io.ReadAll(resp.Response.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		var mangaSearch MangaDexSearch
		if err := json.Unmarshal(body, &mangaSearch); err != nil {
			fmt.Printf("JSON parsing error: %v\nResponse: %s\n", err, string(body))
			return nil, fmt.Errorf("error parsing JSON: %w", err)
		}

		for _, manga := range mangaSearch.Data {
			altTitles := extractAltTitles(manga.Attributes)
			id := manga.ID
			img := extractCoverArtURL(manga.Relationships, p.Api, id)
			format := determineFormat(manga.Type)
			title := extractTitle(manga.Attributes)
			genres := extractGenres(manga.Attributes.Tags)
			description := extractDescription(manga.Attributes.Description)
			countryOfOrigin := extractCountryOfOrigin(manga.Attributes)
			tags := extractTags(manga.Attributes.Tags)
			author := extractAuthor(manga.Relationships)
			publisher := extractPublisher(manga.Relationships)

			results = append(results, types.MediaInfo{
				ID:              id,
				Title:           title,
				Artwork:         nil,
				Synonyms:        altTitles,
				TotalChapters:   helper.ConvertStringToIntPointer(manga.Attributes.LastChapter),
				BannerImage:     nil,
				CoverImage:      &img,
				Color:           nil,
				Year:            &manga.Attributes.Year,
				Status:          &manga.Attributes.Status,
				Genres:          genres,
				Description:     &description,
				Format:          format,
				TotalVolumes:    helper.ConvertStringToIntPointer(manga.Attributes.LastVolume),
				CountryOfOrigin: &countryOfOrigin,
				Tags:            tags,
				Relations:       nil,
				Characters:      nil,
				Author:          &author,
				Publisher:       &publisher,
				Type:            types.TypeManga,
				Rating:          nil,
				Popularity:      nil,
			})
		}
	}

	return results, nil
}

func (p *MangaDexBaseProvider) SearchAdvanced(query string, mediaType types.Type, formats []types.Format, page int, perPage int, genres []string, genresExcluded []string, season types.Season, year int, tags []string, tagsExcluded []string) ([]types.MediaInfo, error) {
	var results []types.MediaInfo

	var genreList []GenreList
	var tagList []GenreList

	if len(tags) > 0 || len(tagsExcluded) > 0 {
		uri, _ := url.Parse(p.Api + "/manga/tag")

		resp, err := p.Request(http.Request{
			URL:    uri,
			Method: "GET",
		}, &p.NeedsProxy)

		if err != nil {
			return nil, err
		}
		defer resp.Response.Body.Close()

		if resp.Response.StatusCode != 200 {
			return nil, fmt.Errorf("unexpected status code: %d", resp.Response.StatusCode)
		}

		if resp.Response.Header.Get("Content-Type") != "application/json" {
			return nil, fmt.Errorf("invalid content type: %s", resp.Response.Header.Get("Content-Type"))
		}

		body, err := io.ReadAll(resp.Response.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		var mangaTags TagResponse
		if err := json.Unmarshal(body, &mangaTags); err != nil {
			fmt.Printf("JSON parsing error: %v\nResponse: %s\n", err, string(body))
			return nil, fmt.Errorf("error parsing JSON: %w", err)
		}

		for _, tag := range mangaTags.Data {
			if tag.Attributes.Group == "theme" {
				genreList = append(genreList, GenreList{
					Name: strings.TrimSpace(tag.Attributes.Name["en"]),
					UID:  tag.ID,
				})
			} else if tag.Attributes.Group == "tag" {
				tagList = append(tagList, GenreList{
					Name: strings.TrimSpace(tag.Attributes.Name["en"]),
					UID:  tag.ID,
				})
			}
		}
	}

	for page := 0; page <= 1; page++ {
		uri, _ := url.Parse(p.Api + "/manga")
		q := uri.Query()

		q.Set("title", query)
		q.Set("limit", "25")
		q.Set("offset", strconv.Itoa(25*page))
		q.Set("order[relevance]", "desc")
		q.Add("contentRating[]", "safe")
		q.Add("contentRating[]", "suggestive")
		q.Add("includes[]", "cover_art")

		if year > 0 {
			q.Set("year", strconv.Itoa(year))
		}
		if len(genres) > 0 {
			var includedTags []string
			for _, genre := range genres {
				for _, item := range genreList {
					if item.Name == genre {
						includedTags = append(includedTags, item.UID)
						break
					}
				}
			}

			if len(includedTags) > 0 {
				q.Set("includedTags[]", strings.Join(includedTags, ","))
				q.Set("includedTagsMode", "AND")
			}
		}
		if len(genresExcluded) > 0 {
			var excludedTags []string
			for _, genre := range genresExcluded {
				for _, item := range genreList {
					if item.Name == genre {
						excludedTags = append(excludedTags, item.UID)
						break
					}
				}
			}

			if len(excludedTags) > 0 {
				q.Set("excludedTags[]", strings.Join(excludedTags, ","))
				if q.Get("includedTagsMode") == "" {
					q.Set("includedTagsMode", "AND")
				}
			}
		}

		if len(tags) > 0 {
			var includedTags []string
			for _, tag := range tags {
				for _, item := range tagList {
					if item.Name == tag {
						includedTags = append(includedTags, item.UID)
						break
					}
				}
			}

			if len(includedTags) > 0 {
				q.Set("includedTags[]", strings.Join(includedTags, ","))
				if q.Get("includedTagsMode") == "" {
					q.Set("includedTagsMode", "AND")
				}
			}
		}

		if len(tagsExcluded) > 0 {
			var excludedTags []string
			for _, tag := range tags {
				for _, item := range tagList {
					if item.Name == tag {
						excludedTags = append(excludedTags, item.UID)
						break
					}
				}
			}

			if len(excludedTags) > 0 {
				q.Set("excludedTags[]", strings.Join(excludedTags, ","))
				if q.Get("includedTagsMode") == "" {
					q.Set("includedTagsMode", "AND")
				}
			}
		}

		uri.RawQuery = q.Encode()

		resp, err := p.Request(http.Request{
			URL:    uri,
			Method: "GET",
		}, &p.NeedsProxy)
		if err != nil {
			return nil, err
		}
		defer resp.Response.Body.Close()

		if resp.Response.StatusCode != 200 {
			return nil, fmt.Errorf("unexpected status code: %d", resp.Response.StatusCode)
		}

		if resp.Response.Header.Get("Content-Type") != "application/json" {
			return nil, fmt.Errorf("invalid content type: %s", resp.Response.Header.Get("Content-Type"))
		}

		body, err := io.ReadAll(resp.Response.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		var mangaSearch MangaDexSearch
		if err := json.Unmarshal(body, &mangaSearch); err != nil {
			fmt.Printf("JSON parsing error: %v\nResponse: %s\n", err, string(body))
			return nil, fmt.Errorf("error parsing JSON: %w", err)
		}

		for _, manga := range mangaSearch.Data {
			altTitles := extractAltTitles(manga.Attributes)
			id := manga.ID
			img := extractCoverArtURL(manga.Relationships, p.Api, id)
			format := determineFormat(manga.Type)
			title := extractTitle(manga.Attributes)
			genres := extractGenres(manga.Attributes.Tags)
			description := extractDescription(manga.Attributes.Description)
			countryOfOrigin := extractCountryOfOrigin(manga.Attributes)
			tags := extractTags(manga.Attributes.Tags)
			author := extractAuthor(manga.Relationships)
			publisher := extractPublisher(manga.Relationships)

			results = append(results, types.MediaInfo{
				ID:              id,
				Title:           title,
				Artwork:         nil,
				Synonyms:        altTitles,
				TotalChapters:   helper.ConvertStringToIntPointer(manga.Attributes.LastChapter),
				BannerImage:     nil,
				CoverImage:      &img,
				Color:           nil,
				Year:            &manga.Attributes.Year,
				Status:          &manga.Attributes.Status,
				Genres:          genres,
				Description:     &description,
				Format:          format,
				TotalVolumes:    helper.ConvertStringToIntPointer(manga.Attributes.LastVolume),
				CountryOfOrigin: &countryOfOrigin,
				Tags:            tags,
				Relations:       nil,
				Characters:      nil,
				Author:          &author,
				Publisher:       &publisher,
				Type:            types.TypeManga,
				Rating:          nil,
				Popularity:      nil,
			})
		}
	}

	return results, nil
}

func (p *MangaDexBaseProvider) GetCurrentSeason() (types.Season, error) {
	return types.SeasonUnknown, nil
}

func (p *MangaDexBaseProvider) GetMedia(id string) (types.MediaInfo, error) {
	uri, _ := url.Parse(p.Api + "/manga/" + id)

	data, err := p.Request(http.Request{
		URL:    uri,
		Method: "GET",
	}, &p.NeedsProxy)
	if err != nil {
		return types.MediaInfo{}, err
	}
	defer data.Response.Body.Close()

	if data.Response.StatusCode != 200 {
		return types.MediaInfo{}, fmt.Errorf("unexpected status code: %d", data.Response.StatusCode)
	}

	if data.Response.Header.Get("Content-Type") != "application/json" {
		return types.MediaInfo{}, fmt.Errorf("invalid content type: %s", data.Response.Header.Get("Content-Type"))
	}

	body, err := io.ReadAll(data.Response.Body)
	if err != nil {
		return types.MediaInfo{}, fmt.Errorf("error reading response body: %w", err)
	}

	var itemData MangaDexItem
	if err := json.Unmarshal(body, &itemData); err != nil {
		fmt.Printf("JSON parsing error: %v\nResponse: %s\n", err, string(body))
		return types.MediaInfo{}, fmt.Errorf("error parsing JSON: %w", err)
	}

	altTitles := extractAltTitles(itemData.Data.Attributes)
	img := extractCoverArtURL(itemData.Data.Relationships, p.Api, id)
	format := determineFormat(itemData.Data.Type)
	title := extractTitle(itemData.Data.Attributes)
	genres := extractGenres(itemData.Data.Attributes.Tags)
	description := extractDescription(itemData.Data.Attributes.Description)
	countryOfOrigin := extractCountryOfOrigin(itemData.Data.Attributes)
	tags := extractTags(itemData.Data.Attributes.Tags)
	author := extractAuthor(itemData.Data.Relationships)
	publisher := extractPublisher(itemData.Data.Relationships)

	return types.MediaInfo{
		ID:              id,
		Title:           title,
		Artwork:         nil,
		Synonyms:        altTitles,
		TotalChapters:   helper.ConvertStringToIntPointer(itemData.Data.Attributes.LastChapter),
		BannerImage:     nil,
		CoverImage:      &img,
		Color:           nil,
		Year:            &itemData.Data.Attributes.Year,
		Status:          &itemData.Data.Attributes.Status,
		Genres:          genres,
		Description:     &description,
		Format:          format,
		TotalVolumes:    helper.ConvertStringToIntPointer(itemData.Data.Attributes.LastVolume),
		CountryOfOrigin: &countryOfOrigin,
		Tags:            tags,
		Relations:       nil,
		Characters:      nil,
		TotalEpisodes:   nil,
		CurrentEpisode:  nil,
		Season:          types.SeasonUnknown,
		Duration:        nil,
		Trailer:         nil,
		Type:            types.TypeManga,
		Rating:          nil,
		Popularity:      nil,
		Author:          &author,
		Publisher:       &publisher,
	}, nil
}

func (p *MangaDexBaseProvider) GetSeasonal(mediaType types.Type, formats []types.Format) (types.SeasonalResponse, error) {
	currentDate := time.Now().AddDate(0, 0, -3)
	createdAtParam := fmt.Sprintf("%04d-%02d-%02dT00:00:00",
		currentDate.Year(),
		int(currentDate.Month()),
		currentDate.Day(),
	)

	trendingUri, _ := url.Parse(p.Api + "/manga?includes[]=cover_art&includes[]=artist&includes[]=author&order[followedCount]=desc&contentRating[]=safe&contentRating[]=suggestive&hasAvailableChapters=true&createdAtSince=" + createdAtParam)
	popularUri, _ := url.Parse(p.Api + "/manga?includes[]=cover_art&includes[]=artist&includes[]=author&order[followedCount]=desc&contentRating[]=safe&contentRating[]=suggestive&hasAvailableChapters=true")
	topUri, _ := url.Parse(p.Api + "/manga?includes[]=cover_art&includes[]=artist&includes[]=author&order[rating]=desc&contentRating[]=safe&contentRating[]=suggestive&hasAvailableChapters=true")
	// https://mangadex.org/titles/seasonal
	seasonalUri, _ := url.Parse(p.Api + "/list/54736a5c-eb7f-4844-971b-80ee171cdf29?includes[]=user")

	trending, err := p.Request(http.Request{
		URL:    trendingUri,
		Method: "GET",
	}, &p.NeedsProxy)
	if err != nil {
		return types.SeasonalResponse{}, err
	}
	defer trending.Response.Body.Close()

	if trending.Response.StatusCode != 200 {
		return types.SeasonalResponse{}, fmt.Errorf("unexpected status code: %d", trending.Response.StatusCode)
	}

	if trending.Response.Header.Get("Content-Type") != "application/json" {
		return types.SeasonalResponse{}, fmt.Errorf("invalid content type: %s", trending.Response.Header.Get("Content-Type"))
	}

	trendingBody, err := io.ReadAll(trending.Response.Body)
	if err != nil {
		return types.SeasonalResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	var trendingData MangaDexSeasonal
	if err := json.Unmarshal(trendingBody, &trendingData); err != nil {
		fmt.Printf("JSON parsing error: %v\nResponse: %s\n", err, string(trendingBody))
		return types.SeasonalResponse{}, fmt.Errorf("error parsing JSON: %w", err)
	}

	popular, err := p.Request(http.Request{
		URL:    popularUri,
		Method: "GET",
	}, &p.NeedsProxy)
	if err != nil {
		return types.SeasonalResponse{}, err
	}
	defer popular.Response.Body.Close()

	if popular.Response.StatusCode != 200 {
		return types.SeasonalResponse{}, fmt.Errorf("unexpected status code: %d", popular.Response.StatusCode)
	}

	if popular.Response.Header.Get("Content-Type") != "application/json" {
		return types.SeasonalResponse{}, fmt.Errorf("invalid content type: %s", popular.Response.Header.Get("Content-Type"))
	}

	popularBody, err := io.ReadAll(popular.Response.Body)
	if err != nil {
		return types.SeasonalResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	var popularData MangaDexSeasonal
	if err := json.Unmarshal(popularBody, &popularData); err != nil {
		fmt.Printf("JSON parsing error: %v\nResponse: %s\n", err, string(popularBody))
		return types.SeasonalResponse{}, fmt.Errorf("error parsing JSON: %w", err)
	}

	top, err := p.Request(http.Request{
		URL:    topUri,
		Method: "GET",
	}, &p.NeedsProxy)
	if err != nil {
		return types.SeasonalResponse{}, err
	}
	defer top.Response.Body.Close()

	if top.Response.StatusCode != 200 {
		return types.SeasonalResponse{}, fmt.Errorf("unexpected status code: %d", top.Response.StatusCode)
	}

	if top.Response.Header.Get("Content-Type") != "application/json" {
		return types.SeasonalResponse{}, fmt.Errorf("invalid content type: %s", top.Response.Header.Get("Content-Type"))
	}

	topBody, err := io.ReadAll(top.Response.Body)
	if err != nil {
		return types.SeasonalResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	var topData MangaDexSeasonal
	if err := json.Unmarshal(topBody, &topData); err != nil {
		fmt.Printf("JSON parsing error: %v\nResponse: %s\n", err, string(topBody))
		return types.SeasonalResponse{}, fmt.Errorf("error parsing JSON: %w", err)
	}

	seasonal, err := p.Request(http.Request{
		URL:    seasonalUri,
		Method: "GET",
	}, &p.NeedsProxy)
	if err != nil {
		return types.SeasonalResponse{}, err
	}
	defer seasonal.Response.Body.Close()

	if seasonal.Response.StatusCode != 200 {
		return types.SeasonalResponse{}, fmt.Errorf("unexpected status code: %d", seasonal.Response.StatusCode)
	}

	if seasonal.Response.Header.Get("Content-Type") != "application/json" {
		return types.SeasonalResponse{}, fmt.Errorf("invalid content type: %s", seasonal.Response.Header.Get("Content-Type"))
	}

	seasonalBody, err := io.ReadAll(seasonal.Response.Body)
	if err != nil {
		return types.SeasonalResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	var seasonalData MangaDexSeasonal
	if err := json.Unmarshal(topBody, &seasonalData); err != nil {
		fmt.Printf("JSON parsing error: %v\nResponse: %s\n", err, string(seasonalBody))
		return types.SeasonalResponse{}, fmt.Errorf("error parsing JSON: %w", err)
	}

	return types.SeasonalResponse{
		Trending: returnFilledManga(p, trendingData),
		Popular:  returnFilledManga(p, popularData),
		Top:      returnFilledManga(p, topData),
		Seasonal: returnFilledManga(p, seasonalData),
	}, nil
}

func (p *MangaDexBaseProvider) GetSchedule() (types.ScheduleResponse, error) {
	return types.ScheduleResponse{}, nil
}

func (p *MangaDexBaseProvider) GetIds() ([]string, error) {
	uri, err := url.Parse("https://raw.githubusercontent.com/ArdaxHz/mangadex-id-map/main/json/manga_map.json")
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	ids, err := p.Request(http.Request{
		URL:    uri,
		Method: "GET",
	}, &p.NeedsProxy)
	if err != nil {
		return nil, err
	}
	defer ids.Response.Body.Close()

	if ids.Response.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", ids.Response.StatusCode)
	}

	idsBody, err := io.ReadAll(ids.Response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var idsData map[string]string
	if err := json.Unmarshal(idsBody, &idsData); err != nil {
		fmt.Printf("JSON parsing error: %v\nResponse: %s\n", err, string(idsBody))
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	var idList []string
	for id := range idsData {
		idList = append(idList, id)
	}

	return idList, nil
}

func returnFilledManga(p *MangaDexBaseProvider, seasonalData MangaDexSeasonal) []types.MediaInfo {
	var results []types.MediaInfo
	for _, data := range seasonalData.Data {
		altTitles := extractAltTitles(data.Attributes)
		img := extractCoverArtURL(data.Relationships, p.Api, data.ID)
		format := determineFormat(data.Type)
		title := extractTitle(data.Attributes)
		genres := extractGenres(data.Attributes.Tags)
		description := extractDescription(data.Attributes.Description)
		countryOfOrigin := extractCountryOfOrigin(data.Attributes)
		tags := extractTags(data.Attributes.Tags)
		author := extractAuthor(data.Relationships)
		publisher := extractPublisher(data.Relationships)

		results = append(results, types.MediaInfo{
			ID:              data.ID,
			Title:           title,
			Artwork:         nil,
			Synonyms:        altTitles,
			TotalChapters:   helper.ConvertStringToIntPointer(data.Attributes.LastChapter),
			BannerImage:     nil,
			CoverImage:      &img,
			Color:           nil,
			Year:            &data.Attributes.Year,
			Status:          &data.Attributes.Status,
			Genres:          genres,
			Description:     &description,
			Format:          format,
			TotalVolumes:    helper.ConvertStringToIntPointer(data.Attributes.LastVolume),
			CountryOfOrigin: &countryOfOrigin,
			Tags:            tags,
			Relations:       nil,
			Characters:      nil,
			Author:          &author,
			Publisher:       &publisher,
			Type:            types.TypeManga,
			Rating:          nil,
			Popularity:      nil,
		})
	}

	return results
}

func extractTitle(attributes ItemAttributes) types.Title {
	title := types.Title{
		English: findTitle(attributes, "en"),
		Romaji:  findRomajiTitle(attributes),
		Native:  findNativeTitle(attributes),
	}
	return title
}

func findTitle(attributes ItemAttributes, lang string) *string {
	for _, altTitle := range attributes.AltTitles {
		if val, exists := altTitle[lang]; exists {
			return &val
		}
	}

	if val, exists := attributes.Title[lang]; exists {
		return &val
	}
	return nil
}

func findRomajiTitle(attributes ItemAttributes) *string {
	if val, exists := attributes.Title["ja-ro"]; exists {
		return &val
	}
	if val, exists := attributes.Title["jp-ro"]; exists {
		return &val
	}

	for _, altTitle := range attributes.AltTitles {
		if val, exists := altTitle["ja-ro"]; exists {
			return &val
		}
		if val, exists := altTitle["jp-ro"]; exists {
			return &val
		}
	}
	return nil
}

func findNativeTitle(attributes ItemAttributes) *string {
	if val, exists := attributes.Title["jp"]; exists {
		return &val
	}
	if val, exists := attributes.Title["ja"]; exists {
		return &val
	}
	if val, exists := attributes.Title["ko"]; exists {
		return &val
	}

	for _, altTitle := range attributes.AltTitles {
		if val, exists := altTitle["jp"]; exists {
			return &val
		}
		if val, exists := altTitle["ja"]; exists {
			return &val
		}
		if val, exists := altTitle["ko"]; exists {
			return &val
		}
	}
	return nil
}

func extractAltTitles(attributes ItemAttributes) []string {
	var altTitles []string
	for _, titleMap := range attributes.AltTitles {
		for _, title := range titleMap {
			altTitles = append(altTitles, title)
		}
	}
	for _, title := range attributes.Title {
		altTitles = append(altTitles, title)
	}
	return altTitles
}

func extractCoverArtURL(relationships []Relationship, apiBase, id string) string {
	for _, rel := range relationships {
		if rel.Type == "cover_art" {
			return fmt.Sprintf("%s/covers/%s/%s.jpg.512.jpg", apiBase, id, rel.ID)
		}
	}
	return ""
}

func extractGenres(tags []Tag) []string {
	var genres []string
	for _, tag := range tags {
		if tag.Attributes.Group == "genre" {
			for _, name := range tag.Attributes.Name {
				genres = append(genres, name)
			}
		}
	}
	return genres
}

func extractDescription(description map[string]string) string {
	if desc, exists := description["en"]; exists {
		return desc
	} else {
		// Return the first description available
		for _, desc := range description {
			return desc
		}

		return "No description available."
	}
}

func extractCountryOfOrigin(attributes ItemAttributes) string {
	if attributes.PublicationDemographic != "" {
		return attributes.PublicationDemographic
	} else {
		return strings.ToUpper(attributes.OriginalLanguage)
	}
}

func extractTags(tags []Tag) []string {
	var tagNames []string
	for _, tag := range tags {
		if tag.Attributes.Group == "theme" {
			if name, exists := tag.Attributes.Name["en"]; exists {
				tagNames = append(tagNames, name)
			}
		}
	}
	return tagNames
}

func extractAuthor(relationships []Relationship) string {
	for _, rel := range relationships {
		if rel.Type == "author" {
			return rel.Attributes.Name
		}
	}
	return ""
}

func extractPublisher(relationships []Relationship) string {
	for _, rel := range relationships {
		if rel.Type == "publisher" {
			return rel.Attributes.Name
		}
	}
	return ""
}

func determineFormat(mangaType string) types.Format {
	if mangaType == "ADAPTATION" {
		return types.FormatManga
	}
	if mangaType == "ONE_SHOT" {
		return types.FormatOneShot
	}
	return types.FormatManga
}

type MangaDexSearch struct {
	Result   string         `json:"result"`
	Response string         `json:"response"`
	Data     []SearchResult `json:"data"`
}

type SearchResult struct {
	ID            string         `json:"id"`
	Attributes    ItemAttributes `json:"attributes"`
	Relationships []Relationship `json:"relationships"`
	Type          string         `json:"type"`
}

type ItemAttributes struct {
	Title                          map[string]string   `json:"title"`
	AltTitles                      []map[string]string `json:"altTitles"`
	Description                    map[string]string   `json:"description"`
	IsLocked                       bool                `json:"isLocked"`
	Links                          map[string]string   `json:"links"`
	OriginalLanguage               string              `json:"originalLanguage"`
	LastVolume                     string              `json:"lastVolume"`
	LastChapter                    string              `json:"lastChapter"`
	PublicationDemographic         string              `json:"publicationDemographic"`
	Status                         string              `json:"status"`
	Year                           int                 `json:"year"`
	ContentRating                  string              `json:"contentRating"`
	Tags                           []Tag               `json:"tags"`
	State                          string              `json:"state"`
	ChapterNumbersResetOnNewVolume bool                `json:"chapterNumbersResetOnNewVolume"`
	CreatedAt                      string              `json:"createdAt"`
	UpdatedAt                      string              `json:"updatedAt"`
	Version                        int                 `json:"version"`
	AvailableTranslatedLanguages   []string            `json:"availableTranslatedLanguages"`
	LatestUploadedChapter          string              `json:"latestUploadedChapter"`
}

type Relationship struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Related    string                 `json:"related"`
	Attributes RelationshipAttributes `json:"attributes"`
}

type Tag struct {
	Attributes TagAttributes `json:"attributes"`
}

type TagAttributes struct {
	Name  map[string]string `json:"name"`
	Group string            `json:"group"`
}

type RelationshipAttributes struct {
	Name string `json:"name"`
}

type GenreList struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
}

type TagResponse struct {
	Result   string     `json:"result"`
	Response string     `json:"response"`
	Data     []TagGroup `json:"data"`
}

type TagGroup struct {
	ID            string             `json:"id"`
	Type          string             `json:"type"`
	Attributes    TagGroupAttributes `json:"attributes"`
	Relationships []Relationship     `json:"relationships"`
}

type TagGroupAttributes struct {
	Name        map[string]string `json:"name"`
	Description map[string]string `json:"description"`
	Group       string            `json:"group"`
	Version     int               `json:"version"`
}

type MangaDexItem struct {
	Result   string   `json:"result"`
	Response string   `json:"response"`
	Data     ItemData `json:"data"`
}

type ItemData struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	Attributes    ItemAttributes `json:"attributes"`
	Relationships []Relationship `json:"relationships"`
}

type MangaDexSeasonal struct {
	Result   string     `json:"result"`
	Response string     `json:"response"`
	Data     []ItemData `json:"data"`
}
