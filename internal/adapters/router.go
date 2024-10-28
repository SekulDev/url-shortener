package adapters

import (
	"github.com/gorilla/mux"
	"net/http"
	"url-shortener/internal/adapters/handlers"
	"url-shortener/internal/infrastructure"
)

func NewRouter(server *infrastructure.Server) *mux.Router {

	router := mux.NewRouter()

	pagesHandler := handlers.NewPagesHandlers(server.Services.TemplateService)
	router.HandleFunc("/", pagesHandler.Index).Methods("GET")

	urlHandlers := handlers.NewUrlHandlers(server.Services.UrlService, server.Services.TemplateService)
	router.HandleFunc("/{short_id}", urlHandlers.ShortIdHandler).Methods("GET")
	router.HandleFunc("/url", urlHandlers.AddUrlHandler).Methods("POST")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	return router
}
