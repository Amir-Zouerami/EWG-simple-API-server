package main

import (
	"net/http"
	"strconv"

	"github.com/Amir-Zouerami/EWG-simple-API-server/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user, err := app.store.Users.GetByID(r.Context(), userID)

	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}

}
