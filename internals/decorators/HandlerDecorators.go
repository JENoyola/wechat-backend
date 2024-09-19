package decorators

import (
	"net/http"
	"wechat-back/internals/database"
)

type HandlerFuncWithDeps func(w http.ResponseWriter, r *http.Request, db *database.DB)

// HandlerDecorator creats a middle man between handlers in order to inject database dependecies
func HandlerDecorator(handler HandlerFuncWithDeps, db *database.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if db == nil {
			db = database.StartDatabase()
		}

		handler(w, r, db)
	}

}
