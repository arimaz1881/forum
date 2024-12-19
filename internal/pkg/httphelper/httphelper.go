package httphelper

import (
	"bytes"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"
	"time"
)

var tmplsCache map[string]*template.Template

func InitTemplates(tmpls map[string]*template.Template) {
	tmplsCache = tmpls
}

type TemplateData struct {
	CurrentYear    int
	IsAuthN        bool
	Role           string
	CanSendRequest bool
	Login		   string
	Form           any
	Errors         map[string]string
}

type File struct {
	FileName    string
	FileReader  multipart.File
	FileSize    int64
	ContentType string
}

type User struct {
	ID             int64
	IsAuthN        bool
	Role           string
	CanSendRequest bool
	Login		   string
}

func GetTmplData(form any, user User) *TemplateData {
	return &TemplateData{
		CurrentYear:    time.Now().Year(),
		IsAuthN:        user.IsAuthN,
		Role:           user.Role,
		CanSendRequest: user.CanSendRequest,
		Login:			user.Login,
		Form:           form,
		Errors:         map[string]string{},
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

func FileFromForm(r *http.Request, key string) (_ *File, err error) {
	files := r.MultipartForm.File[key]
	if len(files) == 0 {
		return nil, nil
	}

	return getFile(files[0])
}

func getFile(fh *multipart.FileHeader) (*File, error) {
	file, err := fh.Open()
	if err != nil {
		return nil, err
	}

	return &File{
		FileName:    fh.Filename,
		FileSize:    fh.Size,
		FileReader:  file,
		ContentType: fh.Header.Get("Content-Type"),
	}, nil
}
