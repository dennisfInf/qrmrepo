package server

import (
	"context"
	"fmt"
	"github.com/enclaive/backend/config"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Server struct {
	echo *echo.Echo
	cfg  config.Config
}

func New(cfg config.Config) *Server {
	e := echo.New()
	e.Server.ReadTimeout = 5 * time.Second
	e.Server.WriteTimeout = 10 * time.Second
	e.Server.IdleTimeout = 120 * time.Second

	return &Server{
		echo: echo.New(),
		cfg:  cfg,
	}
}

func (s *Server) Run() (err error) {
	err = s.registerRoutes()
	if err != nil {
		return
	}

	err = s.echo.Start(fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port))
	if err == http.ErrServerClosed {
		err = nil
	}

	return
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.echo.Shutdown(ctx)
}
