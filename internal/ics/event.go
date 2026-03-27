package ics

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type Calendar struct {
	ProdID string
	Name   string
	Events []Event
}

type Event struct {
	UID         string
	CreatedAt   time.Time
	StartAt     time.Time
	EndAt       time.Time
	Summary     string
	Description string
	URL         string
	ImageURL    string
}

func (c Calendar) Serialize() string {
	var b strings.Builder

	writeLine(&b, "BEGIN:VCALENDAR")
	writeLine(&b, "VERSION:2.0")
	writeLine(&b, "PRODID:"+defaultString(c.ProdID, "-//AniSei//EN"))
	writeLine(&b, "CALSCALE:GREGORIAN")
	writeLine(&b, "METHOD:PUBLISH")

	if c.Name != "" {
		writeLine(&b, "X-WR-CALNAME:"+escapeText(c.Name))
	}

	for _, e := range c.Events {
		serializeEvent(&b, e)
	}

	writeLine(&b, "END:VCALENDAR")

	return b.String()
}

func serializeEvent(b *strings.Builder, e Event) {
	uid := e.UID
	if uid == "" {
		uid = newUID()
	}

	dtstamp := e.CreatedAt
	if dtstamp.IsZero() {
		dtstamp = time.Now().UTC()
	}

	writeLine(b, "BEGIN:VEVENT")
	writeLine(b, "UID:"+uid)
	writeLine(b, "DTSTAMP:"+formatTimeUTC(dtstamp))
	writeLine(b, "DTSTART:"+formatTimeUTC(e.StartAt))
	writeLine(b, "DTEND:"+formatTimeUTC(e.EndAt))

	if e.Summary != "" {
		writeLine(b, "SUMMARY:"+escapeText(e.Summary))
	}
	if e.Description != "" {
		writeLine(b, "DESCRIPTION:"+escapeText(e.Description))
	}
	if e.URL != "" {
		writeLine(b, "URL:"+escapeText(e.URL))
	}
	if e.ImageURL != "" {
		writeLine(b, "ATTACH;FMTTYPE=image/jpeg:"+escapeText(e.URL))
	}
	writeLine(b, "END:VEVENT")
}

func formatTimeUTC(t time.Time) string {
	return t.UTC().Format("20060102T150405Z")
}

func escapeText(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "\r\n", `\n`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, ",", `\,`)
	s = strings.ReplaceAll(s, ";", `\;`)
	return s
}

func writeLine(b *strings.Builder, line string) {
	for _, folded := range foldLine(line) {
		b.WriteString(folded)
		b.WriteString("\r\n")
	}
}

func foldLine(line string) []string {
	const max = 75

	if len(line) <= max {
		return []string{line}
	}

	var out []string
	for len(line) > max {
		out = append(out, line[:max])
		line = " " + line[max:]
		if len(line) <= max {
			break
		}
	}
	out = append(out, line)
	return out
}

func newUID() string {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return fmt.Sprintf("%d@anisei", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf[:]) + "@anisei"
}

func defaultString(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
