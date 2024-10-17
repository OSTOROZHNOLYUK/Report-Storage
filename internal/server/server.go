package server

import (
	"Report-Storage/internal/config"
	"Report-Storage/internal/notifications"
	"Report-Storage/internal/s3cloud"
	"Report-Storage/internal/server/api"
	"Report-Storage/internal/storage/mongodb"
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// Server - структура сервера.
type Server struct {
	srv  *http.Server
	mux  *chi.Mux
	jwt  *jwtauth.JWTAuth
	mail *notifications.SMTP
}

// New - конструктор сервера.
func New(cfg *config.Config) *Server {
	r := chi.NewRouter()
	j := jwtauth.New("HS256", []byte(cfg.JwtSecret), nil, jwt.WithAcceptableSkew(time.Second*30))
	j.ValidateOptions()
	m := notifications.New(cfg.Sender, cfg.SMTPLogin, cfg.SMTPPasswd, cfg.SMTPHost, cfg.SMTPPort)

	server := &Server{
		srv: &http.Server{
			Addr:         cfg.Address,
			Handler:      r,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		mux:  r,
		jwt:  j,
		mail: m,
	}
	return server
}

// Start запускает HTTP сервер в отдельной горутине.
func (s *Server) Start() {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			log.Print("failed to start server")
		}
	}()
}

// API инициализирует все обработчики API.
func (s *Server) API(log *slog.Logger, st *mongodb.Storage, s3 *s3cloud.FileStorage) {
	// Создание заявки.
	s.mux.Post("/api/reports/new", api.AddReport(log, st, s3, s.mail))

	// Безопасные методы.
	s.mux.Post("/api/reports/quad", api.ReportsByPoly(log, st))       // получение заявок в границах многоугольника
	s.mux.Get("/api/reports/all", api.Reports(log, st))               // получение всех заявок
	s.mux.Get("/api/reports/{num}", api.ReportByNum(log, st))         // получение заявки по ее уникальному номеру
	s.mux.Get("/api/reports/filter", api.ReportsWithFilters(log, st)) // получение N заявок с фильтрами
	s.mux.Get("/api/reports/id/{id}", api.ReportByID(log, st))        // получение заявки по ObjectID
	s.mux.Get("/api/reports/radius", api.ReportsByRadius(log, st))    // получение всех заявок в радиусе от заданной точки

	// Методы с проверкой прав.
	s.mux.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(s.jwt))
		r.Use(jwtauth.Authenticator(s.jwt))

		r.Put("/api/reports", api.UpdateReport(log, st, s3, s.mail))                  // обновление всех полей заявки
		r.Patch("/api/reports/status/{num}", api.UpdateStatusReport(log, st, s.mail)) // обновление статуса заявки по ее номеру
		r.Delete("/api/reports/{num}", api.DeleteReport(log, st))                     // удаление заявки по ее номеру
		r.Delete("/api/reports/rejected", api.DeleteRejected(log, st))                // удаление всех заявок со статусом "Отклонена"
		r.Get("/api/reports/statistic", api.Statistic(log, st))                       // получение статистики по всем заявкам
	})
}

// Middleware инициализирует все обработчики middleware.
func (s *Server) Middleware() {
	s.mux.Use(middleware.RequestID)
	s.mux.Use(middleware.RealIP)
	s.mux.Use(middleware.Logger)
	s.mux.Use(middleware.Recoverer)
}

// Shutdown останавливает сервер используя graceful shutdown
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to stop server: %s", err.Error())
	}
}
