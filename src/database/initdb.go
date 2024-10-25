package database

import (
	"context"
	"fmt"
	"os"
)

func CreateTables() {
	anime := `
		CREATE TABLE IF NOT EXISTS anime (
			id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
            slug TEXT,
            "coverImage" TEXT,
            "bannerImage" TEXT,
            trailer TEXT,
            status VARCHAR(255),
            season VARCHAR(255) DEFAULT 'UNKNOWN',
            title JSONB,
            "currentEpisode" REAL,
            mappings JSONB DEFAULT '{}'::JSONB,
            synonyms TEXT[],
            "countryOfOrigin" TEXT,
            description TEXT,
            duration REAL,
            color TEXT,
            year INT,
            rating JSONB,
            popularity JSONB,
            type TEXT,
            format VARCHAR(255) DEFAULT 'UNKNOWN',
            relations JSONB[] DEFAULT '{}'::JSONB[],
            "totalEpisodes" REAL,
            genres TEXT[],
            tags TEXT[],
            episodes JSONB DEFAULT '{"latest": {"updatedAt": 0, "latestEpisode": 0, "latestTitle": ""}, "data": []}'::JSONB,
            "averageRating" REAL,
            "averagePopularity" REAL,
            artwork JSONB[] DEFAULT ARRAY[]::JSONB[],
            characters JSONB[] DEFAULT ARRAY[]::JSONB[]
		);
	`

	manga := `
		CREATE TABLE IF NOT EXISTS manga (
            id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
            slug TEXT,
            "coverImage" TEXT,
            "bannerImage" TEXT,
            status VARCHAR(255),
            title JSONB,
            mappings JSONB DEFAULT '{}'::JSONB,
            synonyms TEXT[],
            "countryOfOrigin" TEXT,
            description TEXT,
            color TEXT,
            year INT,
            rating JSONB,
            popularity JSONB,
            type TEXT,
            format VARCHAR(255) DEFAULT 'UNKNOWN',
            relations JSONB[] DEFAULT '{}'::JSONB[],
            "currentChapter" REAL,
            "totalChapters" REAL,
            "totalVolumes" REAL,
            genres TEXT[],
            tags TEXT[],
            chapters JSONB DEFAULT '{"latest": {"updatedAt": 0, "latestChapter": 0, "latestTitle": ""}, "data": []}'::JSONB,
            "averageRating" REAL,
            "averagePopularity" REAL,
            artwork JSONB[] DEFAULT ARRAY[]::JSONB[],
            characters JSONB[] DEFAULT ARRAY[]::JSONB[]
        );
	`
	extensions := `
		CREATE EXTENSION IF NOT EXISTS pg_trgm;
	`

	functions := `
		create or replace function most_similar(text, text[]) returns double precision
		language sql as $$
			select max(similarity($1,x)) from unnest($2) f(x)
		$$;
	`

	_, err := DB.Exec(context.Background(), anime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create anime table: %v\n", err)
		os.Exit(1)
	}

	_, err = DB.Exec(context.Background(), manga)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create manga table: %v\n", err)
		os.Exit(1)
	}

	_, err = DB.Exec(context.Background(), extensions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create extensions: %v\n", err)
		os.Exit(1)
	}

	_, err = DB.Exec(context.Background(), functions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create functions: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Tables created")
}
