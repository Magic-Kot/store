package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(
	r chi.Router,
	server Server,
	bearerAuth BearerAuth,
) {
	r.Route("/auth/v1", func(r chi.Router) {
		r.Post("/sign-in", handler(server.PostAuthV1SignIn))
		r.Post("/sign-up", handler(server.PostAuthV1SignUp))
		r.Post("/refresh", handler(server.PostAuthV1Refresh))
		r.Post("/logout", handler(server.PostAuthV1Logout))
	})

	r.Route("/user/v1", func(r chi.Router) {
		r.Use(bearerAuth.JWTMiddleware())
		r.Get("/info", handler(server.GetUserV1Info))
	})

	r.Route("/settings/v1", func(r chi.Router) {
		r.Use(bearerAuth.JWTMiddleware())
		r.Patch("/user", handler(server.PatchSettingsV1User))
		//	r.Delete("/user", handler(server.DeleteSettingsV1User))
	})

	r.Route("/bonuses/v1", func(r chi.Router) {
		r.Use(bearerAuth.JWTMiddleware())
		//	r.Post("/friends", handler(server.CreateReferral))
		//	r.Get("/counter", handler(server.CounterReferral))
	})

	//r.Get("/baf/:url", handler(server.GetReferral))
}

func handler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			replyError(r.Context(), w, err)
		}
	}
}
