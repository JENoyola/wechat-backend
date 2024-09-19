package routes

import (
	"wechat-back/internals/decorators"
	"wechat-back/internals/handlers"

	"github.com/go-chi/chi/v5"
)

func UserRoutes(mux *chi.Mux) {

	mux.Post("/nsg", decorators.HandlerDecorator(handlers.NewUserAccountEP, nil))
	mux.Put("/cv", decorators.HandlerDecorator(handlers.UserCodeVerificationEP, nil))
	mux.Get("/sn", decorators.HandlerDecorator(handlers.RequestNewCodeEP, nil))

}
