package database_fetch

import (
	"anify/eltik/go/src/database"
	"anify/eltik/go/src/types"
	"context"

	"github.com/jackc/pgx/v5"
)

func Get(id string, type_ types.Type) (interface{}, error) {
	switch type_ {
	case types.TypeAnime:
		return GetAnimeByID(id)
	case types.TypeManga:
		return GetMangaByID(id)
	}

	return nil, nil
}

// GetAnimeByID fetches an anime by its ID.
func GetAnimeByID(id string) (*types.Anime, error) {
	var anime types.Anime

	err := database.DB.QueryRow(context.Background(), `
		SELECT id, artwork, "averagePopularity", "averageRating", "bannerImage", characters, color,
		       "countryOfOrigin", "coverImage", "currentEpisode", description, duration, episodes,
		       format, genres, mappings, popularity, rating, relations, season, slug, status,
		       synonyms, tags, title, "totalEpisodes", trailer, type, year
		FROM anime
		WHERE id = $1
	`, id).Scan(&anime.ID, &anime.Artwork, &anime.AveragePopularity, &anime.AverageRating, &anime.BannerImage, &anime.Characters, &anime.Color, &anime.CountryOfOrigin, &anime.CoverImage, &anime.CurrentEpisode, &anime.Description, &anime.Duration, &anime.Episodes, &anime.Format, &anime.Genres, &anime.Mappings, &anime.Popularity, &anime.Rating, &anime.Relations, &anime.Season, &anime.Slug, &anime.Status, &anime.Synonyms, &anime.Tags, &anime.Title, &anime.TotalEpisodes, &anime.Trailer, &anime.Type, &anime.Year)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &anime, nil
}

func GetMangaByID(id string) (*types.Manga, error) {
	var manga types.Manga

	err := database.DB.QueryRow(context.Background(), `
		SELECT id, artwork, "averagePopularity", "averageRating", "bannerImage", color, "countryOfOrigin",
			   "coverImage", "currentChapter", description, format, genres, mappings, popularity,
			   rating, relations, slug, status, synonyms, title, "totalChapters",
			   "totalVolumes", type, year							
		FROM manga
		WHERE id = $1
	`, id).Scan(&manga.ID, &manga.Artwork, &manga.AveragePopularity, &manga.AverageRating, &manga.BannerImage, &manga.Color, &manga.CountryOfOrigin, &manga.CoverImage, &manga.CurrentChapter, &manga.Description, &manga.Format, &manga.Genres, &manga.Mappings, &manga.Popularity, &manga.Rating, &manga.Relations, &manga.Slug, &manga.Status, &manga.Synonyms, &manga.Title, &manga.TotalChapters, &manga.TotalVolumes, &manga.Type, &manga.Year)
	if err != nil {
		if err == pgx.ErrNoRows {
			println("No rows found!")
			return &types.Manga{}, nil
		}
		return nil, err
	}

	return &manga, nil
}
