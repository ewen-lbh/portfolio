package pages

import (
	"encoding/json"
)

// dump dumps the given content into a string by marshalling it to JSON.
func dump(v any) string {
	raw, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(raw)
}
