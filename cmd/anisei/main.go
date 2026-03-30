package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/naoki914/anisei/internal/anilist"
	"github.com/naoki914/anisei/internal/ics"
)

// const endpoint = "https://graphql.anilist.co"

func main() {
	// Handle user set variables
	var (
		seasonFlag = flag.String("season", "", "Specify the season (winter/summer/spring/fall)")
		year       = flag.Int("year", time.Now().Year(), "Specify the year")
		output     = flag.String("output", "calendar.ics", "Output calendar file name")
	)

	flag.Parse()
	season := anilist.GetSeasonFromString(*seasonFlag)

	//fmt.Printf("%s-%s | %d | %s\n", *seasonFlag, season, *year, *output)
	// Verify year
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	getSeasonAnimeSchedule(ctx, season, *year, *output)
}

func getSeasonAnimeSchedule(ctx context.Context, season anilist.MediaSeason, year int, filename string) {
	c := anilist.New(nil)
	//SeasonData, err := c.SearchSeasonAnime(ctx, anilist.SeasonWinter, 2026)
	SeasonData, err := c.GetAllSeasonAnimeSchedule(ctx, season, year)
	if err != nil {
		panic(err)
	}

	cal := buildAnimeCalendar(SeasonData)
	err = WriteFile(filename, cal)
	if err != nil {
		panic(err)
	}
	//	fmt.Printf("%v\n\n", cal)

}

func buildAnimeCalendar(animeList []anilist.Media) ics.Calendar {
	cal := ics.Calendar{
		ProdID: "WinterCalendar",
		Name:   "Winter2016",
	}

	for _, anime := range animeList {
		//fmt.Printf("%v\n", anime)
		for _, ep := range anime.AiringSchedule.Nodes {
			start := time.Unix(ep.AiringAt, 0).UTC()
			end := start.Add(30 * time.Minute)

			cal.Events = append(cal.Events, ics.Event{
				UID:         fmt.Sprintf("anilist-%s", slug(anime.Title.Romaji)),
				CreatedAt:   time.Now().UTC(),
				StartAt:     start,
				EndAt:       end,
				Summary:     fmt.Sprintf("%s - %d", anime.Title.Romaji, ep.Episode),
				Description: fmt.Sprintf("%s episode %d airing", anime.Title.Romaji, ep.Episode),
				URL:         anime.SiteURL,
				ImageURL:    anime.CoverImage.ExtraLarge,
			})
		}
	}

	return cal
}

func slug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
func WriteFile(path string, cal ics.Calendar) error {
	fmt.Println(path)
	return os.WriteFile(path, []byte(cal.Serialize()), 0644)
}
