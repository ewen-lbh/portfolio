package shared

import (
	"context"
	"io"

	"github.com/a-h/templ"
	ortfodb "github.com/ortfo/db"
)

func HTML(html ortfodb.HTMLString) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, string(html))
		return
	})
}

func CSSDeclaration(property string, value string) string {
	return property + ": " + value + ";\n"
}

func CSS(declarations map[string]map[string]string) templ.Component {
	var css string
	for selector, decls := range declarations {
		css += selector + " {\n"
		for property, value := range decls {
			css += "\t" + CSSDeclaration(property, value)
		}
		css += "}\n"
	}
	return HTML(ortfodb.HTMLString("<style>\n" + css + "</style>"))
}
