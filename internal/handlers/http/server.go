package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
	"wbL0/internal/app/config"
	"wbL0/internal/usecase"
)

type Server struct {
	config     *config.Config
	logger     *zap.SugaredLogger
	experiment usecase.Experiment
}

func NewServer(config config.Config, logger zap.SugaredLogger, expUC usecase.Experiment) *Server {
	return &Server{
		config:     &config,
		logger:     &logger,
		experiment: expUC,
	}
}

func (s *Server) Run() error {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(cors.New())
	s.RegisterRoutes(app)
	return app.Listen(":8080")
}
func (s *Server) RegisterRoutes(app *fiber.App) {
	app.Get("/api/v1/health-check", s.HealthCheck)
	app.Post("/api/v1/pullOrder", s.PullOrder)
	app.Post("/api/v1/getOrderById", s.getOrder)
}
