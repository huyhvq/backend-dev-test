package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.NotFound(app.notFound)
	mux.MethodNotAllowed(app.methodNotAllowed)
	mux.Use(app.recoverPanic)

	mux.Get("/_/status", app.status)

	mux.Route("/v1", func(r chi.Router) {
		r.Get("/organizations/search", app.organizationSearch)
		r.Post("/hubs", app.createHub)
		r.Post("/teams", app.createTeam)
		r.Patch("/teams/{id}/join-hub", app.joinIntoHub)
		r.Post("/users", app.createUser)
		r.Patch("/users/{id}/join-team", app.joinIntoTeam)
	})

	return mux
}
