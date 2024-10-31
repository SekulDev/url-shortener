package service

import (
	"html/template"
	"io"
)

type TemplateService struct {
	templates       *template.Template
	recaptchaPublic string
}

func NewTemplateService(recaptchaPublic string) *TemplateService {

	templates := template.Must(template.ParseGlob("web/tmpl/*.gohtml"))
	//template.Must(templates.ParseGlob("web/tmpl/partials/*.gohtml"))

	return &TemplateService{
		templates:       templates,
		recaptchaPublic: recaptchaPublic,
	}
}

func (t *TemplateService) RenderPage(w io.Writer, name string, data interface{}) error {
	err := t.Render(w, name, data)
	if err != nil {
		return err
	}

	err2 := t.templates.ExecuteTemplate(w, "layout.gohtml", map[string]interface{}{
		"RecaptchaPublic": t.recaptchaPublic,
	})
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
