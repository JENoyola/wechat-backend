package routes

import (
	"wechat-back/internals/decorators"
	"wechat-back/internals/handlers"

	"github.com/go-chi/chi/v5"
)

func ChatRoutes(mux *chi.Mux) {

	mux.Handle("/uchat", decorators.HandlerDecorator(handlers.HandleP2PConnectionEP, nil))
	mux.Handle("/gchat", decorators.HandlerDecorator(handlers.HandleGroupConnectionsEP, nil))

}
