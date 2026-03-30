# AniSei

 A Go application that generates calendar files (__*.ics__) for anime airing schedules from AniList.

## Features

 - Fetch seasonal anime schedules from AniList GraphQL API
 - Generate calendar events automatically from anime airing data
 - Support for seasonal anime schedules (Winter/Summer/Spring/Fall)

## Quick Start

### Clone and build

```
git clone github.com/naoki914/anisei
cd anisei
go build -o anisei ./cmd/anisei
```

### Usage

Usage

Run Anisei via command-line flags:

./anisei -season winter -year 2026 -output calendar.ics

Command-Line Flags

 | Flag | Description                  | Default Value     |
 |------|------------------------------|-------------------|
 | -season | Season (winter/spring/summer/fall)    | Current Season |
 | -year   | Year to schedule                   | Current year      |
 | -output | Output ICS calendar file name     | calendar.ics    |

#### Examples

 # Winter 2026 calendar
 ./anisei -season winter -year 2026 -output 2026-winter.ics

 # Spring 2025 calendar
 ./anisei -season spring -year 2025 -output 2025-spring.ics

 # Custom output name
 ./anisei -season fall -year 2026 -output my-animelist.ics

### Configuration

 Customize these aspects in cmd/anisei/main.go:
 - Season selection (anilist.SeasonWinter|Spring|Summer|Fall)
 - Year of the anime schedule
 - Output file path
 - Event duration per episode

### Requirements

 - Go 1.20+
 - Access to AniList GraphQL API

### Output

 Generates an .ics file compatible with most calendar applications (Outlook, Google Calendar, Apple Calendar,
 Thunderbird, etc.).
