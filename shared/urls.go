package shared

import (
	"github.com/a-h/templ"
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

func Asset(path string) templ.SafeURL {
	if path == "" {
		panic("asset path cannot be empty")
	}
	if IsDev() {
		return templ.URL("http://localhost:8079/" + path)
	}
	return templ.URL("https://assets.ewen.works/" + path)
}
