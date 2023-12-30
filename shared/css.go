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

// Luminance returns the luminance of the given color, normalized from 0-1.
func Luminance(hexstring string) float64 {
	if hexstring == "black" {
		return 0
	}

	if hexstring == "white" {
		return 1
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
	return (0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)) / 255
}

// IsColorDark returns true if the color is dark, false otherwise.
func IsColorDark(hexstring string) bool {
	return Luminance(hexstring) < 0.5
}

// Contrast returns the contrast ratio between two colors.
func Contrast(hexstring1 string, hexstring2 string) float64 {
	l1 := Luminance(hexstring1)
	l2 := Luminance(hexstring2)
	if l1 > l2 {
		return (l1 + 0.05) / (l2 + 0.05)
	}
	return (l2 + 0.05) / (l1 + 0.05)
}

// ReadableOn returns a color that is readable on the given color.
func ReadableOn(color string) string {
	if color == "" {
		return ""
	}
	if IsColorDark(Color(color)) {
		return "#fff"
	}
	return "#000"
}

// Color ensures that the color given has an octothorpe (#) prefix.
func Color(color string) string {
	if color == "" {
		return ""
	}
	if _, err := hex.DecodeString(color); err == nil {
		return "#" + color
	}
	return color
}
