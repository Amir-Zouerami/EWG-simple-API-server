package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Amir-Zouerami/EWG-simple-API-server/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1, //TODO: change after auth
	}

	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post, err := app.store.Posts.GetByID(r.Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}
