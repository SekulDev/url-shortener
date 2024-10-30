package handlers

import (
	"net/http"
	"url-shortener/internal/app/service"
)

type PagesHandlers struct {
	templateService *service.TemplateService
}

func NewPagesHandlers(templateService *service.TemplateService) *PagesHandlers {
	return &PagesHandlers{
		templateService: templateService,
	}
}

func (h *PagesHandlers) Index(w http.ResponseWriter, r *http.Request) {
	err := h.templateService.RenderPage(w, "index.gohtml", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
