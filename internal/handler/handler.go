package handler

import (
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/sprint-id/catsocialx/internal/cfg"
	"github.com/sprint-id/catsocialx/internal/service"
)

type Handler struct {
	router  *chi.Mux
	service *service.Service
	cfg     *cfg.Cfg
}

func NewHandler(router *chi.Mux, service *service.Service, cfg *cfg.Cfg) *Handler {
	handler := &Handler{router, service, cfg}
	handler.registRoute()

	return handler
}

func (h *Handler) registRoute() {

	r := h.router
	var tokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte(h.cfg.JWTSecret), nil, jwt.WithAcceptableSkew(30*time.Second))

	userH := newUserHandler(h.service.User)
	catH := newCatHandler(h.service.Cat)
	matchH := newMatchHandler(h.service.Match)

	r.Use(middleware.RedirectSlashes)

	r.Post("/v1/user/register", userH.Register)
	r.Post("/v1/user/login", userH.Login)

	// protected route
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Patch("/v1/user", userH.UpdateAccount)

		r.Post("/v1/cat", catH.AddCat)
		r.Get("/v1/cat", catH.GetCat)
		r.Get("/v1/cat/{id}", catH.GetCatByID)
		r.Put("/v1/cat/{id}", catH.UpdateCat)
		r.Delete("/v1/cat/{id}", catH.DeleteCat)

		r.Post("/v1/cat/match", matchH.MatchCat)
		r.Get("/v1/cat/match", matchH.GetMatch)

		r.Post("/v1/cat/match/approve", matchH.ApproveMatch)
		r.Post("/v1/cat/match/reject", matchH.RejectMatch)
		r.Delete("/v1/cat/match/{id}", matchH.DeleteMatch)
	})
}
