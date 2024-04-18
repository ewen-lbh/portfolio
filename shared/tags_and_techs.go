package shared

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/metal3d/go-slugify"
	ortfodb "github.com/ortfo/db"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// stringsLooselyMatch checks if s1 is equal to any of sn, but case-insensitively.
func stringsLooselyMatch(s1 string, sn ...string) bool {
	collator := collate.New(language.English, collate.Loose)
	if len(sn) == 0 {
		return false
	} else {
		return collator.CompareString(s1, sn[0]) == 0 || stringsLooselyMatch(s1, sn[1:]...)
	}
}

// String returns the string representation of the external site.
// Should be the one used in URLs, as GetDistFilepath uses this.
func (s Site) String() string {
	return s.Name
}

// String returns the string representation of the tag.
// Should be the one used in URLs, as GetDistFilepath uses this.
func (t Tag) String() string {
	return t.URLName()
}

// String returns the string representation of the technology.
// Should be the one used in URLs, as GetDistFilepath uses this.
func (t Technology) String() string {
	return t.Name
}

// URLName computes the identifier to use in the tag's page's URL
func (t Tag) URLName() string {
	return slugify.Marshal(t.Plural)
}

// ReferredToBy returns whether the given name refers to the tag
func (t *Tag) ReferredToBy(name string) bool {
	return stringsLooselyMatch(name, t.Plural, t.Singular, t.URLName()) || stringsLooselyMatch(name, t.Aliases...)
}

// ReferredToBy returns whether the given name refers to the tech
func (t *Technology) ReferredToBy(name string) bool {
	return stringsLooselyMatch(name, t.Slug, t.Name) || stringsLooselyMatch(name, t.Aliases...)
}

func (t Technology) ExtendsTech(tech Technology) bool {
	for _, extendedTech := range t.Extends {
		if tech.ReferredToBy(extendedTech) {
			return true
		}
	}
	return false
}

// CalculateTimeSpent updates the time spent on the technology, using wakatime data
func (t *Technology) CalculateTimeSpent(techs []Technology) (time.Duration, error) {
	if len(timeSpentOnTechs) > 0 {
		totalDuration := time.Duration(0)
		for techName, duration := range timeSpentOnTechs {
			if t.ReferredToBy(techName) {
				totalDuration += duration
				extendors := make([]Technology, 0)
				for _, otherTech := range techs {
					if otherTech.ExtendsTech(*t) {
						extendors = append(extendors, otherTech)
					}
				}
				if len(extendors) > 0 {
					for _, otherTech := range extendors {
						timeSpentOnOtherTech, err := otherTech.CalculateTimeSpent(techs)
						if err != nil {
							return 0, fmt.Errorf("while resolving time spent of tech %s that extends %s: %w", otherTech.Name, t.Name, err)
						}
						totalDuration += timeSpentOnOtherTech
					}
				}
			}
		}
		t.TimeSpent = totalDuration
		return totalDuration, nil
	}

	var stats struct {
		Data wakatimeUserStats `json:"data"`
	}

	resp, err := wakatimeRequest("users/current/stats/all_time")
	if err != nil {
		return 0, fmt.Errorf("while fetcing user stats: %w", err)
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&stats)
	if err != nil {
		return 0, fmt.Errorf("while decoding user stats: %w", err)
	}

	for _, tech := range stats.Data.Languages {
		timeSpentOnTechs[tech.Name] = time.Duration(tech.TotalSeconds) * time.Second
	}

	return t.CalculateTimeSpent(techs)
}

func TagsOf(all []Tag, metadata ortfodb.WorkMetadata) []Tag {
	tags := make([]Tag, len(metadata.Tags))
	for i, t := range metadata.Tags {
		tags[i] = LookupTag(all, t)
	}
	return tags
}

func TechsOf(all []Technology, metadata ortfodb.WorkMetadata) []Technology {
	techs := make([]Technology, len(metadata.MadeWith))
	for i, t := range metadata.MadeWith {
		techs[i] = LookupTech(all, t)
	}
	return techs
}

func LookupTag(all []Tag, name string) Tag {
	for _, t := range all {
		if t.ReferredToBy(name) {
			return t
		}
	}
	panic(fmt.Errorf("no tag found with name %q", name))
}

func LookupTech(all []Technology, name string) Technology {
	for _, t := range all {
		if t.ReferredToBy(name) {
			return t
		}
	}
	panic(fmt.Errorf("no technology found with name %q", name))
}
