package server

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/SHSanderland/EffMobTest/pkg/config"
	"github.com/SHSanderland/EffMobTest/pkg/handlers"
	"github.com/SHSanderland/EffMobTest/pkg/storage"
	"github.com/go-chi/chi/v5"
)

func InitServer(l *slog.Logger, cfg *config.Config, db storage.Storage) {
	const fn = "server.InitServer"
	log := l.With(
		slog.String("fn", fn),
		slog.String("Address server", cfg.Addr),
	)

	srv := http.Server{
		Addr:         cfg.Addr,
		Handler:      initMux(l, db),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info("Start server!")

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Error("error to start server", slog.String("err", err.Error()))

		panic(err)
	}
}

func initMux(log *slog.Logger, db storage.Storage) *chi.Mux {
	router := chi.NewRouter()
	h := handlers.InitHandlers(log, db)

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/subscriptions", h.CreateSubscription)
	})

	return router
}
