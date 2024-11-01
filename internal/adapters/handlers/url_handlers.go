package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"url-shortener/internal/adapters/middleware"
	"url-shortener/internal/app/service"
	urlS "url-shortener/internal/app/service/url"
)

type UrlHandlers struct {
	urlService      urlS.UrlService
	templateService service.TemplateService
}

func NewUrlHandlers(urlService urlS.UrlService, templateService service.TemplateService) *UrlHandlers {
	return &UrlHandlers{
		urlService:      urlService,
		templateService: templateService,
	}
}

func (h *UrlHandlers) ShortIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortId := vars["short_id"]

	url, err := h.urlService.ResolveShortUrl(shortId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.LongUrl, http.StatusMovedPermanently)
}

func (h *UrlHandlers) AddUrlHandler(w http.ResponseWriter, r *http.Request) {
	ip, ok := r.Context().Value(middleware.IpContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url := r.FormValue("longurl")
	recaptchaToken := r.FormValue("g-recaptcha-response")

	result, err := h.urlService.AddUrl(url, ip, recaptchaToken)
	if err != nil {
		err := h.templateService.Render(w, "add_url_response.gohtml", map[string]interface{}{
			"error": err.Error(),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		err := h.templateService.Render(w, "add_url_response.gohtml", map[string]interface{}{
			"short_url": result.ShortId,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
