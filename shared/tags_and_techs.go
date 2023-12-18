package shared

import (
	"strings"

	"github.com/metal3d/go-slugify"
)

// stringsLooselyMatch checks if s1 is equal to any of sn, but case-insensitively.
func stringsLooselyMatch(s1 string, sn ...string) bool {
	for _, s2 := range sn {
		if strings.EqualFold(s1, s2) {
			return true
		}
	}
	return false
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
