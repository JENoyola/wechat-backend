package routes

import (
	"wechat-back/internals/decorators"
	"wechat-back/internals/handlers"

	"github.com/go-chi/chi/v5"
)

func UserRoutes(mux *chi.Mux) {

	mux.Post("/nsg", decorators.HandlerWProvidersDecorator(handlers.NewUserAccountEP, nil, nil))
	mux.Put("/cv", decorators.HandlerDecorator(handlers.UserCodeVerificationEP, nil))
	mux.Get("/sn", decorators.HandlerDecorator(handlers.RequestNewCodeEP, nil))
	mux.Post("/l", decorators.HandlerDecorator(handlers.LoginEP, nil))
	mux.Get("/ulkup", decorators.HandlerDecorator(handlers.SearchUsersEP, nil))
}
