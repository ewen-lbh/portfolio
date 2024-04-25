package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ewen-lbh/portfolio/shared"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

var md = goldmark.New(goldmark.WithExtensions(
	extension.Table,
	extension.Footnote,
	extension.CJK,
	extension.DefinitionList,
	extension.Strikethrough,
	extension.Typographer,
	extension.TaskList,
	&frontmatter.Extender{},
	highlighting.Highlighting,
))

func loadBlogEntries(searchIn string) (entries []shared.BlogEntry, err error) {
	files, err := os.ReadDir(searchIn)
	if err != nil {
		return entries, fmt.Errorf("while reading directory %q: %w", searchIn, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		contents, err := os.ReadFile(filepath.Join(searchIn, file.Name()))
		if err != nil {
			return entries, fmt.Errorf("while reading %q: %w", file.Name(), err)
		}

		var out bytes.Buffer
		ctx := parser.NewContext()
		md.Convert(contents, &out, parser.WithContext(ctx))

		var entry shared.BlogEntry
		err = frontmatter.Get(ctx).Decode(&entry)
		if err != nil {
			return entries, fmt.Errorf("while decoding frontmatter of %q: %w", file.Name(), err)
		}

		entry.Slug = strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		entry.Content = out.String()

		entries = append(entries, entry)
	}

	return entries, nil
}
