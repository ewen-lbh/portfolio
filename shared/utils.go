package shared

import (
	"encoding/hex"
	"sort"
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

func Color(color string) string {
	if color == "" {
		return "inherit"
	}
	if _, err := hex.DecodeString(color); err == nil {
		return "#" + color
	}
	return color
}
