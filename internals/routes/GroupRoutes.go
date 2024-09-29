package routes

import (
	"net/http"
	"wechat-back/internals/decorators"
	"wechat-back/internals/handlers"

	"github.com/go-chi/chi/v5"
)

func GroupHandlers(mux *chi.Mux) http.Handler {

	mux.Post("/cg", decorators.HandlerWProvidersDecorator(handlers.CreateNewGroupEP, nil, nil))
	mux.Put("/ugi", decorators.HandlerDecorator(handlers.UpdateGroupInfoEP, nil))
	mux.Put("/uga", decorators.HandlerDecorator(handlers.UpdateGroupAdminsEP, nil))
	mux.Put("/igp", decorators.HandlerDecorator(handlers.UpdateGroupParticipantsEP, nil))
	mux.Delete("/dgp", decorators.HandlerDecorator(handlers.DeleteGroupEP, nil))
	mux.Get("/sgp", decorators.HandlerDecorator(handlers.SearchGroupsEP, nil))

	return mux
}
