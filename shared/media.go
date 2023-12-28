package shared

import (
	"fmt"
	"strings"

	ortfodb "github.com/ortfo/db"
)

func MediaSrcSet(media ortfodb.Media) string {
	sizes := make([]string, 0, len(media.Thumbnails))
	for size, path := range media.Thumbnails {
		sizes = append(sizes, fmt.Sprintf("%s %dw", Media(path), size))
	}
	return strings.Join(sizes, ", ")
}
