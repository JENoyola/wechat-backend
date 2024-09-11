package routes

import (
	"wechat-back/internals/handlers"

	"github.com/go-chi/chi/v5"
)

func HealtRoutes(mux *chi.Mux) {

	mux.Get("/", handlers.ServerHealthCheckEP)

}
