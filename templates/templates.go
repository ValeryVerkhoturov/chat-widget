package templates

import (
	"embed"
	"github.com/ValeryVerkhoturov/chat/utils/requestUtils"
	"html/template"
)

var (
	//go:embed all:templates/*
	TemplateFS embed.FS

	// HTML parsed templates
	HTML *template.Template
)

func init() {
	var err error
	HTML, err = requestUtils.TemplateParseFSRecursive(TemplateFS, ".html", true, nil)
	if err != nil {
		panic(err)
	}
}
