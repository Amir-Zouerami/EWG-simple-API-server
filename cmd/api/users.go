package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Amir-Zouerami/EWG-simple-API-server/internal/store"
	"github.com/go-chi/chi/v5"
)

type userContextKey string

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

const (
	userCtxKey userContextKey = "user"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	toBeFollowedUser := getUserFromContext(r)

	var currUser FollowUser // TODO: needs to be removed when authentication is added
	if err := readJSON(w, r, &currUser); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Follow(ctx, toBeFollowedUser.ID, currUser.UserID); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, toBeFollowedUser); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	toBeUnfollowedUser := getUserFromContext(r)

	var currUser FollowUser // TODO: needs to be removed when authentication is added
	if err := readJSON(w, r, &currUser); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Unfollow(ctx, toBeUnfollowedUser.ID, currUser.UserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, toBeUnfollowedUser); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) userContextMIddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)

		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, userID)

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

		ctx = context.WithValue(ctx, userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtxKey).(*store.User)

	return user
}
