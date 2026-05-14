package handler

import (
	"html/template"
	"net/http"

	"github.com/shamhub/exercise/pkg/errorlib"
	"github.com/shamhub/exercise/pkg/httpservice"
)

func GetUserId(ctx *httpservice.RequestContextForTemplate) (interface{}, error) {
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil, errorlib.NewResponseError(http.StatusInternalServerError, err.Error())
	}

	data := struct {
		Name        string
		Visitor     string
		Description string
		Socials     map[string]string
	}{
		Name:        "Alice",
		Visitor:     "Visitor",
		Description: "Welcome to website",
		Socials: map[string]string{
			"Twitter/X": "@sdf",
			"LinkedIn":  "fdfgdfg",
			"Instagram": "example",
		},
	}

	return httpservice.TemplateData{
		TemplateHandle: tmpl,
		Data:           data,
		TemplateName:   "home.html",
	}, nil
}
