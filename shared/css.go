package shared

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"

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

func OnFocus(class templ.CSSClass, rules Declarations) templ.Component {
	selector := fmt.Sprintf(".%s:focus", class.ClassName())
	return CSS(Selectors{
		selector: rules,
	})
}

// IsColorDark returns true if the color is dark, false otherwise.
func IsColorDark(hexstring string) bool {
	if hexstring == "black" {
		return true
	}

	if hexstring == "white" {
		return false
	}

	// Convert (#)rrggbb to decimal
	rgb, err := strconv.ParseInt(strings.TrimPrefix(hexstring, "#"), 16, 32)
	if err != nil {
		panic(fmt.Sprintf("while checking if %s is a dark color: %s", hexstring, err))
	}

	// Extract components
	r := rgb >> 16
	g := (rgb >> 8) & 0xFF
	b := rgb & 0xFF

	// Calculate luminance (see ITU-R BT.709), normalized from 0-255 to 0-1
	luminance := (0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)) / 255

	return luminance < 0.5
}

func ReadableOn(color string) string {
	if color == "" {
		return ""
	}
	if IsColorDark(Color(color)) {
		return "#fff"
	}
	return "#000"
}

func Color(color string) string {
	if color == "" {
		return ""
	}
	if _, err := hex.DecodeString(color); err == nil {
		return "#" + color
	}
	return color
}
