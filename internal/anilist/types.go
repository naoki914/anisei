package anilist

type MediaSeason string

const (
	SeasonWinter MediaSeason = "WINTER"
	SeasonSpring MediaSeason = "SPRING"
	SeasonSummer MediaSeason = "SUMMER"
	SeasonFall   MediaSeason = "FALL"
)

// Title represents anime title translations
type Title struct {
	Romaji  string `json:"romaji"`
	English string `json:"english,omitempty"`
	Native  string `json:"native,omitempty"`
}

// NextAiringEpisode represents upcoming episode info
type NextAiringEpisode struct {
	Episode         int   `json:"episode,omitempty"`
	AiringAt        int64 `json:"airingAt,omitempty"`
	TimeUntilAiring int   `json:"timeUntilAiring,omitempty"`
}

// Media represents an anime/media item from the API
type Media struct {
	ID             int                `json:"id"`
	Title          Title              `json:"title"`
	Season         string             `json:"season"`
	SeasonYear     int                `json:"seasonYear"`
	Status         string             `json:"status"`
	Format         string             `json:"format"`
	Episodes       *int               `json:"episodes,omitempty"`
	SiteURL        string             `json:"siteUrl"`
	NextAiring     *NextAiringEpisode `json:"nextAiringEpisode,omitempty"`
	CurrentEpisode int                `json:"episodesWatched,omitempty"`
	CoverImage     CoverImage         `json:"coverImage"`
	Genres         []string           `json:"genres"`
	AiringSchedule struct {
		Nodes []AiringSchedule `json:"nodes"`
	} `json:"airingSchedule"`
}

type PageInfo struct {
	Total       int  `json:"total"`
	CurrentPage int  `json:"currentPage"`
	LastPage    int  `json:"lastPage"`
	HasNextPage bool `json:"hasNextPage"`
	PerPage     int  `json:"perPage"`
}

type GQLError struct {
	Message string `json:"message"`
}

// Response holds the complete GraphQL response
type Response[T any] struct {
	Data   T          `json:"data"`
	Errors []GQLError `json:"errors,omitempty"`
}

type AiringSchedule struct {
	Episode   int   `json:"episode"`
	AiringAt  int64 `json:"airingAt"`
	TimeUntil int   `json:"timeUntilAiring"`
}

type CoverImage struct {
	Medium     string `json:"medium,omitempty"`
	Large      string `json:"large,omitempty"`
	ExtraLarge string `json:"extraLarge.omitempty"`
}
type SeasonAnimeData struct {
	Page struct {
		PageInfo PageInfo `json:"pageInfo"`
		Media    []Media  `json:"media"`
	} `json:"Page"`
}

type AnimeScheduleData struct {
	Media struct {
		ID             int   `json:"id"`
		Title          Title `json:"title"`
		AiringSchedule struct {
			Nodes []AiringSchedule `json:"nodes"`
		} `json:"airingSchedule"`
	} `json:"Media"`
}
