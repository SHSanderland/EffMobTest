// Пакет server используется для инициализации сервера
// и роутера приложения.
package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/SHSanderland/EffMobTest/pkg/config"
	"github.com/SHSanderland/EffMobTest/pkg/handlers"
	"github.com/SHSanderland/EffMobTest/pkg/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/SHSanderland/EffMobTest/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// InitServer Инициализация сервера.
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

	log.Info("Start serve../r!")

	go gracefulShutdown(log, &srv)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Error("error to start server", slog.String("err", err.Error()))

		panic(err)
	}

	log.Info("Closing database...")

	db.CloseConnection()
}

// initMux Инициализация роутера.
func initMux(log *slog.Logger, db storage.Storage) *chi.Mux {
	router := chi.NewRouter()
	h := handlers.InitHandlers(log, db)

	router.Use(
		middleware.Recoverer,
		middleware.RequestID,
		middleware.Logger,
	)

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/subscriptions", h.CreateSubscription)
		r.Get("/subscriptions/{id}", h.ReadSubscription)
		r.Put("/subscriptions/{id}", h.UpdateSubscription)
		r.Delete("/subscriptions/{id}", h.DeleteSubscription)
		r.Get("/subscriptions", h.ListSubscription)
		r.Get("/subscriptions/cost", h.CostSubscription)
	})

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	return router
}

// gracefulShutdown Функция для постепенного выключения сервера.
// Слушает сигналы ОС. Запускать в горутине.
func gracefulShutdown(log *slog.Logger, srv *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	log.Info("Received shutdown signal, gracefully shutting down...")

	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Warn("failed to shutdown server", slog.String("err", err.Error()))
	}
}
