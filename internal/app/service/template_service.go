package service

import (
	"html/template"
	"io"
)

type TemplateService struct {
	templates *template.Template
}

func NewTemplateService() *TemplateService {

	templates := template.Must(template.ParseGlob("web/tmpl/*.html"))
	template.Must(templates.ParseGlob("web/tmpl/partials/*.html"))

	return &TemplateService{
		templates: templates,
	}
}

func (t *TemplateService) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
