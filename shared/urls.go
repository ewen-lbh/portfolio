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

func Asset(path string) string {
	if path == "" {
		panic("asset path cannot be empty")
	}
	if IsDev() {
		return "http://localhost:8079/" + path
	}
	return "https://assets.ewen.works/" + path
}
