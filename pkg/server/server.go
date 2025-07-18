package server

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/SHSanderland/EffMobTest/pkg/config"
)

func InitServer(l *slog.Logger, cfg *config.Config) {
	const fn = "server.InitServer"
	log := l.With(
		slog.String("fn", fn),
		slog.String("Address server", cfg.Addr),
	)

	srv := http.Server{
		Addr:         cfg.Addr,
		Handler:      nil,
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
