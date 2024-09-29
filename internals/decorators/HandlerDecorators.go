package decorators

import (
	"net/http"
	"wechat-back/internals/database"
	"wechat-back/providers/media"
)

type handlerFuncWithDeps func(w http.ResponseWriter, r *http.Request, db database.DBHUB)
type handlerWithProviders func(w http.ResponseWriter, r *http.Request, db database.DBHUB, m media.MediaHUB)

// HandlerDecorator creats a middle man between handlers in order to inject database dependecies
func HandlerDecorator(handler handlerFuncWithDeps, db database.DBHUB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if db == nil {
			db = database.StartDatabase()
		}

		handler(w, r, db)
	}

}

func HandlerWProvidersDecorator(handler handlerWithProviders, db database.DBHUB, m media.MediaHUB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if db == nil {
			db = database.StartDatabase()
		}

		if m == nil {
			m, _ = media.NewMediaService()

		}

		handler(w, r, db, m)

	}
}
