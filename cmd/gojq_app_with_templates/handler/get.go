package handler

import (
	"html/template"
	"net/http"
	"strings"
	"time"

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
		Features    []string
	}{
		Name:        "Alice",
		Visitor:     "Visitor",
		Description: "Welcome to website",
		Socials: map[string]string{
			"Twitter/X": "@sdf",
			"LinkedIn":  "fdfgdfg",
			"Instagram": "example",
		},
		Features: []string{
			"Customizable Products",
			"24/7 Customer Support",
			"Reliable and Secure",
		},
	}

	return httpservice.TemplateData{
		TemplateHandle: tmpl,
		Data:           data,
		TemplateName:   "home.html",
	}, nil
}

func toUpper(str string) string {
	return strings.ToUpper(str)
}

func formatDate(t time.Time) string {
	return t.Format("January 2, 2007")
}

func GetUser(ctx *httpservice.RequestContextForTemplate) (interface{}, error) {

	funcMap := template.FuncMap{
		"toUpper":     toUpper,
		"formateDate": formatDate,
	}

	tmpl, err := template.New("functions.html").Funcs(funcMap).ParseFiles("templates/functions.html")
	if err != nil {
		return nil, errorlib.NewResponseError(http.StatusInternalServerError, err.Error())
	}

	data := struct {
		Name        string
		CurrentDate time.Time
		Number      int
		Items       []string
	}{
		Name:        "John Doe",
		CurrentDate: time.Now(),
		Number:      7,
		Items:       []string{"Apples", "Oranges", "Bananas"},
	}
	return httpservice.TemplateData{
		TemplateHandle: tmpl,
		Data:           data,
		TemplateName:   "functions.html",
	}, nil
}
