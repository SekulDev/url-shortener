package adapters

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/{short_id}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods("GET")

	router.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {

	}).Methods("POST")

	return router
}
