package shared

import (
	"context"
	"fmt"
	"io"

	"github.com/a-h/templ"
	ortfodb "github.com/ortfo/db"
)

type Declarations map[string]string
type Selectors map[string]Declarations

func HTML(html ortfodb.HTMLString) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, string(html))
		return
	})
}

func CSSDeclaration(property string, value string) string {
	return property + ": " + value + ";\n"
}

func CSS(declarations Selectors) templ.Component {
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

func OnHover(class templ.CSSClass, rules Declarations) templ.Component {
	selector := fmt.Sprintf(".%s:hover, .%s:focus-visible", class.ClassName(), class.ClassName())
	return CSS(Selectors{
		selector: rules,
	})
}
