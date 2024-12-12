package httphelper

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var tmplsCache map[string]*template.Template

func InitTemplates(tmpls map[string]*template.Template) {
	tmplsCache = tmpls
}

type TemplateData struct {
	CurrentYear int
	IsAuthN     bool
	Form        any
	Errors      map[string]string
}

func GetTmplData(form any, isAuthN bool) *TemplateData {
	return &TemplateData{
		CurrentYear: time.Now().Year(),
		IsAuthN:     isAuthN,
		Form:        form,
		Errors:      map[string]string{},
	}
}

func Render(w http.ResponseWriter, status int, page string, data *TemplateData) {
	ts, ok := tmplsCache[page]
	if !ok {
		// TODO: parse error
		fmt.Println("tmplsCache[page]", tmplsCache[page])
		return
	}

	buf := &bytes.Buffer{}

	if err := ts.ExecuteTemplate(buf, "base", data); err != nil {
		fmt.Println(err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}
