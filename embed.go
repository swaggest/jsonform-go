package jsonform

import (
	"embed"
	"html/template"

	"github.com/vearutop/statigz"
)

var (
	// FS holds embedded static assets.
	//
	//go:embed static/*
	staticAssets embed.FS

	staticServer = statigz.FileServer(staticAssets, statigz.FSPrefix("static"))
)

func loadTemplate(fileName string) *template.Template {
	tpl, err := staticAssets.ReadFile("static/" + fileName)
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New("htmlResponse").Parse(string(tpl))
	if err != nil {
		panic(err)
	}

	return tmpl
}
