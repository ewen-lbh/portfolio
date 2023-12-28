package shared

import (
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/goodsign/monday"
)

type MapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

func SortedKeys[V any](m map[string]V) []MapEntry[string, V] {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	entries := make([]MapEntry[string, V], len(keys))
	for i, k := range keys {
		entries[i] = MapEntry[string, V]{k, m[k]}
	}
	return entries
}

func Merge[K comparable, V any](m1 map[K]V, m2 map[K]V) map[K]V {
	out := make(map[K]V)
	for k, v := range m1 {
		out[k] = v
	}
	for k, v := range m2 {
		out[k] = v
	}
	return out
}

func Keys[K comparable, V any](m map[K]V) (keys []K) {
	for k := range m {
		keys = append(keys, k)
	}
	return
}

func MondayLocale(language string) monday.Locale {
	switch language {
	case "fr":
		return monday.LocaleFrFR
	default:
		return monday.LocaleEnUS
	}
}

func FormatDate(date time.Time, format string, locale string) string {
	if date.Year() == 9999 {
		return ""
	}
	return monday.Format(date, format, MondayLocale(locale))
}

func DomainOfURL(urlString string) string {
	url, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}

	return strings.TrimPrefix(url.Hostname(), "www.")
}
