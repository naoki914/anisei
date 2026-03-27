package anilist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const Endpoint = "https://graphql.anilist.co"

type Request struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type Client struct {
	http *http.Client
}

func New(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{http: httpClient}
}

func (c *Client) doRaw(ctx context.Context, r Request) ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, Endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "AniSei/0.1 (anilist client)")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// include body to make debugging painless
		return nil, fmt.Errorf("anilist http error: %s: %s", resp.Status, string(raw))
	}
	return raw, nil
}

func decodeResponse[T any](raw []byte) (T, error) {
	var zero T
	// fmt.Println(string(raw))
	var out Response[T]

	if err := json.Unmarshal(raw, &out); err != nil {
		fmt.Println(1)
		return zero, fmt.Errorf("decode response : %w", err)
	}

	if len(out.Errors) > 0 {
		return zero, fmt.Errorf("anilist graphql error: %s", out.Errors[0].Message)
	}

	return out.Data, nil
}

func doTyped[T any](ctx context.Context, c *Client, r Request) (T, error) {
	raw, err := c.doRaw(ctx, r)
	if err != nil {
		var zero T
		return zero, err
	}
	return decodeResponse[T](raw)
}

const searchSeasonAnimeQuery = `
query ($season: MediaSeason!, $seasonYear: Int!, $page: Int!, $perPage: Int!) {
  Page(page: $page, perPage: $perPage) {
    pageInfo {
      total
      currentPage
      lastPage
      hasNextPage
      perPage
    }
    media(
      type: ANIME
      season: $season
      seasonYear: $seasonYear
      sort: POPULARITY_DESC
    ) {
      id
      title {
        romaji
        english
        native
      }
      season
      seasonYear
      status
      format
      episodes
      siteUrl
      nextAiringEpisode {
        episode
        airingAt
        timeUntilAiring
      }
			coverImage {
				extraLarge
				large
			}
			genres
    }
  }
}
`

func (c *Client) SearchSeasonAnime(ctx context.Context, season MediaSeason, seasonYear int) ([]Media, error) {
	var perPage int = 50
	var pageIx int = 1
	var media []Media
	var pageInfo PageInfo
	var firstRun bool = true

	for {
		req := Request{
			Query: searchSeasonAnimeQuery,
			Variables: map[string]any{
				"season":     season,
				"seasonYear": seasonYear,
				"page":       pageIx,
				"perPage":    perPage,
			},
		}
		raw, err := c.doRaw(ctx, req)
		if err != nil {
			return nil, err
		}
		data, err := decodeResponse[SeasonAnimeData](raw)
		if err != nil {
			return nil, err
		}
		if firstRun {
			media = data.Page.Media
			pageInfo = data.Page.PageInfo
			firstRun = false
		} else {
			media = append(media, data.Page.Media...)
			pageInfo.CurrentPage = data.Page.PageInfo.CurrentPage
			pageInfo.HasNextPage = data.Page.PageInfo.HasNextPage
			pageInfo.LastPage = data.Page.PageInfo.LastPage
			pageInfo.Total += data.Page.PageInfo.Total
			pageInfo.PerPage = data.Page.PageInfo.PerPage
		}
		if !pageInfo.HasNextPage {
			break
		}
		pageIx++
	}
	return media, nil
}

const animeScheduleQuery = `
query ($id: Int!) {
  Media(id: $id, type: ANIME) {
    id
    title {
      romaji
      english
      native
    }
    airingSchedule(perPage: 50) {
      nodes {
        episode
        airingAt
        timeUntilAiring
      }
    }
  }
}
`

func (c *Client) GetAnimeSchedule(ctx context.Context, id int) (AnimeScheduleData, error) {
	req := Request{
		Query: animeScheduleQuery,
		Variables: map[string]any{
			"id": id,
		},
	}

	raw, err := c.doRaw(ctx, req)
	if err != nil {
		var zero AnimeScheduleData
		return zero, err
	}

	//	fmt.Println(string(raw))
	return decodeResponse[AnimeScheduleData](raw)

}

const searchAllSeasonAnimeQuery = `
query ($season: MediaSeason!, $seasonYear: Int!, $page: Int!, $perPage: Int!) {
  Page(page: $page, perPage: $perPage) {
    pageInfo {
      total
      currentPage
      lastPage
      hasNextPage
      perPage
    }
    media(
      type: ANIME
      season: $season
      seasonYear: $seasonYear
      sort: POPULARITY_DESC
    ) {
      id
      title {
        romaji
        english
        native
      }
      season
      seasonYear
      status
      format
      episodes
      siteUrl
      nextAiringEpisode {
        episode
        airingAt
        timeUntilAiring
      }
			coverImage {
				extraLarge
				large
			}
			genres
			airingSchedule {
				nodes {
					airingAt
					episode
					id
					timeUntilAiring
				}
			}
    }
  }
}
`

func (c *Client) GetAllSeasonAnimeSchedule(ctx context.Context, season MediaSeason, seasonYear int) ([]Media, error) {
	perPage := 50
	pageIx := 1
	firstRun := true

	var media []Media
	for {
		req := Request{
			Query: searchAllSeasonAnimeQuery,
			Variables: map[string]any{
				"season":     season,
				"seasonYear": seasonYear,
				"page":       pageIx,
				"perPage":    perPage,
			},
		}
		raw, err := c.doRaw(ctx, req)
		if err != nil {
			return nil, err
		}
		data, err := decodeResponse[SeasonAnimeData](raw)
		if err != nil {
			return nil, err
		}
		if firstRun {
			media = data.Page.Media
			firstRun = false
		} else {
			media = append(media, data.Page.Media...)
		}
		if !data.Page.PageInfo.HasNextPage {
			break
		}
		pageIx++
	}

	//	fmt.Println(string(raw))
	return media, nil

}
