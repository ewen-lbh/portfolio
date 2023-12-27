package shared

import (
	"time"

	ortfodb "github.com/ortfo/db"
)

type LocalizableString map[string]string

func (ls LocalizableString) Localized(locale string) string {
	return ls[locale]
}

type Site struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Purpose  string `yaml:"purpose"`
	Username string `yaml:"username"`
}

type Technology struct {
	Slug         string             `yaml:"slug"`
	Name         string             `yaml:"name"`
	By           string             `yaml:"by"`
	Files        []string           `yaml:"files"`
	Aliases      []string           `yaml:"aliases"`
	LearnMoreURL string             `yaml:"learn more at"`
	Description  ortfodb.HTMLString `yaml:"description"`
	Extends      []string           `yaml:"extends"`
	TimeSpent    time.Duration
}

func (tech Technology) Works(db ortfodb.Database) (worksWithTech []ortfodb.AnalyzedWork) {
	for _, work := range db.Works() {
		for _, name := range work.Metadata.MadeWith {
			if tech.ReferredToBy(name) {
				worksWithTech = append(worksWithTech, work)
			}
		}
	}
	return
}

type Tag struct {
	Singular     string             `yaml:"singular"`
	Plural       string             `yaml:"plural"`
	Aliases      []string           `yaml:"aliases"`
	Description  ortfodb.HTMLString `yaml:"description"`
	LearnMoreURL string             `yaml:"learn more at"`
}

func (tag Tag) Works(db ortfodb.Database) (worksWithTag []ortfodb.AnalyzedWork) {
	for _, work := range db.Works() {
		for _, tagName := range work.Metadata.Tags {
			if tag.ReferredToBy(tagName) {
				worksWithTag = append(worksWithTag, work)
			}
		}
	}
	return
}

type Collection struct {
	ID           string
	Title        map[string]string `yaml:"title"`
	Aliases      []string          `yaml:"aliases"`
	Includes     string            `yaml:"includes"`
	Description  map[string]string `yaml:"description"`
	LearnMoreURL string            `yaml:"learn more at"`
}

type Collections map[string]Collection

func (cs Collections) URLsToNames(canonicalOnly bool, locale string) map[string]string {
	urlsToNames := make(map[string]string)
	for id, collection := range cs {
		if !canonicalOnly {
			for _, alias := range collection.Aliases {
				urlsToNames[alias] = collection.Title[locale]
			}
		}
		urlsToNames[id] = collection.Title[locale]
	}
	return urlsToNames
}

func (cs Collections) ThatIncludeWork(work ortfodb.AnalyzedWork, workIDs []string, tags []Tag, techs []Technology) []Collection {
	out := make([]Collection, 0)
	for _, c := range cs {
		if ok, err := c.Contains(work, workIDs, tags, techs); ok || err != nil {
			if err != nil {
				panic(err)
			}

			out = append(out, c)
		}
	}
	return out
}
