package routes

import (
	"wechat-back/internals/decorators"
	"wechat-back/internals/handlers"

	"github.com/go-chi/chi/v5"
)

func ChatRoutes(mux *chi.Mux) {

	mux.Handle("/uchat", decorators.HandlerWProvidersDecorator(handlers.HandleP2PConnectionEP, nil, nil))
	mux.Handle("/gchat", decorators.HandlerWProvidersDecorator(handlers.HandleGroupConnectionsEP, nil, nil))

}
