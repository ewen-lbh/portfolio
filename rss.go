package main

import (
	"github.com/ewen-lbh/portfolio/shared"
	"github.com/gorilla/feeds"
)

func BlogRssFeed(entries []shared.BlogEntry) (*feeds.Feed, error) {
	feed := &feeds.Feed{
		Title:       "ewen.works blog",
		Link:        &feeds.Link{Href: "https://ewen.works/blog"},
		Description: "Blog of Ewen",
		Author:      &feeds.Author{Name: "Ewen Le Bihan", Email: "hey@ewen.works"},
		Created:     earliestEntry(entries...).Date,
	}

	for _, entry := range entries {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       entry.Title,
			Link:        &feeds.Link{Href: "https://ewen.works/blog/" + entry.Slug},
			Description: entry.Content,
			Author:      feed.Author,
			Created:     entry.Date,
		})
	}

	return feed, nil
}

func earliestEntry(entries ...shared.BlogEntry) *shared.BlogEntry {
	if len(entries) == 0 {
		return &shared.BlogEntry{}
	}

	if len(entries) == 1 {
		return &entries[0]
	}

	a, b := entries[0], earliestEntry(entries[1:]...)
	if a.Date.Before(b.Date) {
		return &a
	}
	return b
}
