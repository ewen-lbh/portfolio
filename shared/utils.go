package shared

import "sort"

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
