// Package tmplfunc contains extension functions for templates
package tmplfunc

import (
	"html/template"
	"io/fs"
	"net/url"

	"github.com/bmatcuk/doublestar/v4"
)

func FuncMap(
	h HashedNamers, templates fs.FS,
	partialRenderer func(name string, data any) (template.HTML, error),
) template.FuncMap {
	return template.FuncMap{
		"queryEscape":  url.QueryEscape,
		"appHashed":    h.AppHashed,
		"staticHashed": h.StaticHashed,
		"templates": func(pattern string) []string {
			results, err := doublestar.Glob(templates, pattern, doublestar.WithFilesOnly())
			if err != nil {
				return nil
			}
			return results
		},
		"partial": partialRenderer,
	}
}
