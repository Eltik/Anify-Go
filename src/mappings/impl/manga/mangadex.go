package manga

import (
	"anify/eltik/go/src/lib/impl/request"
	"anify/eltik/go/src/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type MangaDexProvider struct {
	types.BaseMangaProvider
	Api string
}

func NewMangaDexProvider() *MangaDexProvider {
	return &MangaDexProvider{
		BaseMangaProvider: types.BaseMangaProvider{
			RateLimit:          250,
			Id:                 "mangadex",
			Url:                "https://mangadex.org",
			Formats:            []types.Format{types.FormatManga},
			ProviderType:       types.ProviderTypeManga,
			NeedsProxy:         true,
			UseGoogleTranslate: false,
		},
		Api: "https://api.mangadex.org",
	}
}

func (p *MangaDexProvider) Search(query string, format types.Format, year int) ([]types.Result, error) {
	var results []types.Result

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

			results = append(results, types.Result{
				ID:         id,
				Title:      title,
				AltTitles:  altTitles,
				Year:       manga.Attributes.Year,
				Format:     format,
				Img:        &img,
				ProviderId: p.Id,
			})
		}
	}

	if len(results) > 0 {
		fmt.Println(results[0].Title)
	}

	return results, nil
}

func extractAltTitles(attributes Attributes) []string {
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

func determineFormat(mangaType string) types.Format {
	if mangaType == "ADAPTATION" {
		return types.FormatManga
	}
	if mangaType == "ONE_SHOT" {
		return types.FormatOneShot
	}
	return types.FormatManga
}

func extractTitle(attributes Attributes) string {
	titleKeys := []string{"en", "ja-ro", "jp-ro", "jp", "ja", "ko"}
	for _, key := range titleKeys {
		if title, exists := attributes.Title[key]; exists {
			return title
		}
		for _, altTitle := range attributes.AltTitles {
			if title, exists := altTitle[key]; exists {
				return title
			}
		}
	}
	return ""
}

// FetchChapters fetches chapters for a specific manga by ID.
func (p *MangaDexProvider) FetchChapters(id string) ([]types.Chapter, error) {
	return nil, nil
}

// FetchRecent fetches the most recent manga.
func (p *MangaDexProvider) FetchRecent() ([]types.Manga, error) {
	return nil, nil
}

// FetchPages fetches pages for a given chapter.
func (p *MangaDexProvider) FetchPages(id string, proxy bool, chapter *types.Chapter) (interface{}, error) {
	return nil, nil
}

// ProxyCheck checks if the provider can access the API through a proxy.
func (p *MangaDexProvider) ProxyCheck() (bool, error) {
	return false, nil
}

// PadNum pads a number with leading zeros.
func (p *MangaDexProvider) PadNum(number string, places int) string {
	return p.BaseMangaProvider.PadNum(number, places)
}

func (p *MangaDexProvider) Request(config http.Request, proxyRequest *bool) (request.Response, error) {
	return p.BaseMangaProvider.Request(config, proxyRequest)
}

type MangaDexSearch struct {
	Result   string         `json:"result"`
	Response string         `json:"response"`
	Data     []SearchResult `json:"data"`
}

type SearchResult struct {
	ID            string         `json:"id"`
	Attributes    Attributes     `json:"attributes"`
	Relationships []Relationship `json:"relationships"`
	Type          string         `json:"type"`
}

type Attributes struct {
	AltTitles []map[string]string `json:"altTitles"`
	Title     map[string]string   `json:"title"`
	Year      int                 `json:"year"`
}

type Relationship struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Related    string                 `json:"related"`
	Attributes map[string]interface{} `json:"attributes"`
}
