package types

type Type string

type ProviderType string

const (
	TypeAnime Type = "ANIME"
	TypeManga Type = "MANGA"
)

const (
	ProviderTypeAnime ProviderType = "ANIME"
	ProviderTypeManga ProviderType = "MANGA"
)

type Format string

type Provider struct {
	ID  string
	URL string
}

type Result struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	AltTitles  []string `json:"altTitles"`
	Year       int      `json:"year"`
	Format     Format   `json:"format"`
	Img        *string  `json:"img,omitempty"`
	ProviderId string   `json:"providerId"`
}

const (
	FormatTV      Format = "TV"
	FormatTVShort Format = "TV_SHORT"
	FormatMovie   Format = "MOVIE"
	FormatSpecial Format = "SPECIAL"
	FormatOVA     Format = "OVA"
	FormatONA     Format = "ONA"
	FormatMusic   Format = "MUSIC"
	FormatManga   Format = "MANGA"
	FormatNovel   Format = "NOVEL"
	FormatOneShot Format = "ONE_SHOT"
	FormatUnknown Format = "UNKNOWN"
)

type Rating map[string]float64
type Popularity map[string]float64

type Season string

const (
	SeasonWinter  Season = "WINTER"
	SeasonSpring  Season = "SPRING"
	SeasonSummer  Season = "SUMMER"
	SeasonFall    Season = "FALL"
	SeasonUnknown Season = "UNKNOWN"
)

type Title struct {
	Romaji  *string
	English *string
	Native  *string
}

type Mapping struct {
	ID           string
	ProviderID   string
	Similarity   float64
	ProviderType *string
}

type Artwork struct {
	Type       string
	Img        string
	ProviderID string
}

type Character struct {
	Name       string
	Image      string
	VoiceActor VoiceActor
}

type VoiceActor struct {
	Name  string
	Image string
}

type Relations struct {
	ID           string
	Type         Type
	Title        Title
	Format       Format
	RelationType string
}

type Episode struct {
	ID          string
	Title       string
	Number      int
	IsFiller    bool
	Img         *string
	HasDub      bool
	Description *string
	Rating      *float64
	UpdatedAt   *int64
}

type EpisodeData struct {
	ProviderID string
	Episodes   []Episode
}

type Chapter struct {
	ID        string
	Title     string
	Number    int
	Rating    *float64
	UpdatedAt *int64
	Mixdrop   *string
}

type ChapterData struct {
	ProviderID string
	Chapters   []Chapter
}

type Status string

const (
	StatusFinished    Status = "FINISHED"
	StatusReleasing   Status = "RELEASING"
	StatusNotYetAired Status = "NOT_YET_RELEASED"
	StatusCancelled   Status = "CANCELLED"
	StatusHiatus      Status = "HIATUS"
	StatusUnknown     Status = "UNKNOWN"
)

type Media struct {
	ID                string
	Slug              string
	CoverImage        *string
	BannerImage       *string
	Trailer           *string
	Status            *Status
	Season            Season
	Title             Title
	CurrentEpisode    *int
	Mappings          []Mapping
	Synonyms          []string
	CountryOfOrigin   *string
	Description       *string
	Duration          *int
	Color             *string
	Year              *int
	Rating            Rating
	Popularity        Popularity
	AverageRating     *float64
	AveragePopularity *float64
	Type              Type
	Genres            []string
	Format            Format
	Relations         []Relations
	TotalEpisodes     *int
	Episodes          EpisodeCollection
	Tags              []string
	Artwork           []Artwork
	Characters        []Character

	CurrentChapter *int
	TotalVolumes   *int
	Publisher      *string
	Author         *string
	TotalChapters  *int
	Chapters       ChapterCollection
}

type Anime struct {
	ID                string
	Slug              string
	CoverImage        *string
	BannerImage       *string
	Trailer           *string
	Status            *Status
	Season            Season
	Title             Title
	CurrentEpisode    *int
	Mappings          []Mapping
	Synonyms          []string
	CountryOfOrigin   *string
	Description       *string
	Duration          *int
	Color             *string
	Year              *int
	Rating            Rating
	Popularity        Popularity
	AverageRating     *float64
	AveragePopularity *float64
	Type              Type
	Genres            []string
	Format            Format
	Relations         []Relations
	TotalEpisodes     *int
	Episodes          EpisodeCollection
	Tags              []string
	Artwork           []Artwork
	Characters        []Character
}

type Manga struct {
	ID                string
	Slug              string
	CoverImage        *string
	BannerImage       *string
	Status            *Status
	Title             Title
	Mappings          []Mapping
	Synonyms          []string
	CountryOfOrigin   *string
	Description       *string
	CurrentChapter    *int
	TotalVolumes      *int
	Color             *string
	Year              *int
	Rating            Rating
	Popularity        Popularity
	AverageRating     *float64
	AveragePopularity *float64
	Genres            []string
	Type              Type
	Format            Format
	Relations         []Relations
	Publisher         *string
	Author            *string
	TotalChapters     *int
	Chapters          ChapterCollection
	Tags              []string
	Artwork           []Artwork
	Characters        []Character
}

type EpisodeCollection struct {
	Latest struct {
		UpdatedAt     int64
		LatestEpisode int
		LatestTitle   string
	}
	Data []EpisodeData
}

type ChapterCollection struct {
	Latest struct {
		UpdatedAt     int64
		LatestChapter int
		LatestTitle   string
	}
	Data []ChapterData
}

type MediaInfo struct {
	ID              string
	Title           Title
	Artwork         []Artwork
	Synonyms        []string
	TotalEpisodes   *int
	CurrentEpisode  *int
	BannerImage     *string
	CoverImage      *string
	Color           *string
	Season          Season
	Year            *int
	Status          *string
	Genres          []string
	Description     *string
	Format          Format
	Duration        *int
	Trailer         *string
	CountryOfOrigin *string
	Tags            []string
	Relations       []Relations
	Characters      []Character
	Type            Type
	Rating          *float64
	Popularity      *float64
	TotalChapters   *int
	TotalVolumes    *int
	Author          *string
	Publisher       *string
}
