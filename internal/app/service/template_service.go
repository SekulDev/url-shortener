package service

import (
	"html/template"
	"io"
)

type TemplateService struct {
	templates *template.Template
}

func NewTemplateService() *TemplateService {

	templates := template.Must(template.ParseGlob("web/tmpl/*.gohtml"))
	//template.Must(templates.ParseGlob("web/tmpl/partials/*.gohtml"))

	return &TemplateService{
		templates: templates,
	}
}

func (t *TemplateService) RenderPage(w io.Writer, name string, data interface{}) error {
	err := t.Render(w, name, data)
	if err != nil {
		return err
	}

	err2 := t.templates.ExecuteTemplate(w, "layout.gohtml", data)
	if err2 != nil {
		return err2
	}
	return nil
}

func (t *TemplateService) Render(w io.Writer, name string, data interface{}) error {
	err := t.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		return err
	}
	return nil
}
