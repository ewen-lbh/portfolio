package shared

import (
	ortfodb "github.com/ortfo/db"
)

func Media(path ortfodb.FilePathInsideMediaRoot) string {
	if string(path) == "" {
		return ""
	}
	var origin string
	if IsDev() {
		origin = "http://localhost:8080"
	} else {
		origin = "https://media.ewen.works"
	}
	return origin + "/" + string(path)
}
