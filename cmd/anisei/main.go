package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/naoki914/anisei/internal/anilist"
	"github.com/naoki914/anisei/internal/ics"
)

// const endpoint = "https://graphql.anilist.co"

func main() {
	c := anilist.New(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//SeasonData, err := c.SearchSeasonAnime(ctx, anilist.SeasonWinter, 2026)
	data, err := c.GetAllSeasonAnimeSchedule(ctx, anilist.SeasonWinter, 2026)
	if err != nil {
		panic(err)
	}
	fmt.Println("=============================")
	for ix, v := range data {

		fmt.Printf("%d - %v\n", ix, v)
		fmt.Println("-----------------------------")
	}
	/**
	for _, v := range SeasonData {
		//fmt.Printf("- %s [%s]: %d | %s\n", v.Title.English, v.Title.Romaji, v.ID, v.CoverImage.Large)
		//anime, err := c.GetAnimeSchedule(ctx, v.ID)
		 err != nil {
			panic(err)
		}
		cal := BuildAnimeCalendar(anime.Media.Title.Romaji, v.SiteURL, v.CoverImage.ExtraLarge, anime.Media.AiringSchedule.Nodes)
		WriteFile("./Test1.ics", cal)
		fmt.Printf("%v\n\n", cal)
		break
	}*/
}

func BuildAnimeCalendar(title string, url string, imageUrl string, episodes []anilist.AiringSchedule) ics.Calendar {
	cal := ics.Calendar{
		ProdID: "-//AniSei//EN",
		Name:   title,
	}

	for _, ep := range episodes {
		start := time.Unix(ep.AiringAt, 0).UTC()
		end := start.Add(30 * time.Minute)

		cal.Events = append(cal.Events, ics.Event{
			UID:         fmt.Sprintf("anilist-%s-ep-%d@anisei", slug(title), ep.Episode),
			CreatedAt:   time.Now().UTC(),
			StartAt:     start,
			EndAt:       end,
			Summary:     fmt.Sprintf("%s - Episode %d", title, ep.Episode),
			Description: fmt.Sprintf("%s episode %d airing", title, ep.Episode),
			URL:         url,
			ImageURL:    imageUrl,
		})
	}

	return cal
}

func slug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
func WriteFile(path string, cal ics.Calendar) error {
	return os.WriteFile(path, []byte(cal.Serialize()), 0644)
}
