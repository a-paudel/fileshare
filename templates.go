package main

import (
	"embed"
	"io"
	"io/fs"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
)

//go:embed templates/*
var embeddedFiles embed.FS

type EchoTemplates struct {
	templates map[string]*template.Template
}

func (t *EchoTemplates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates[name].ExecuteTemplate(w, "base.html", data)
}

func RegisterTemplates(app *echo.Echo) {
	templateMap := make(map[string]*template.Template)

	subdir, err := fs.Sub(embeddedFiles, "templates")
	if err != nil {
		panic(err)
	}

	err = fs.WalkDir(subdir, ".", func(path string, d fs.DirEntry, err error) error {
		// if path not html skip
		if !strings.HasSuffix(path, ".html") {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if path == "base.html" {
			return nil
		}
		name := strings.TrimSuffix(path, ".html")

		templateMap[name] = template.Must(template.ParseFS(subdir, path, "base.html"))
		return nil
	})

	if err != nil {
		panic(err)
	}

	app.Renderer = &EchoTemplates{
		templates: templateMap,
	}

}
