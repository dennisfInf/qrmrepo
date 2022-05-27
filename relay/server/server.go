package server

import (
	"context"
	"fmt"
	"github.com/enclaive/relay/config"
	"github.com/enclaive/relay/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type LookupRepository interface {
	GetByUsername(ctx context.Context, username string) (models.Lookup, error)
	Set(ctx context.Context, l models.Lookup) error
}

type RepositoryManager interface {
	IsEmptyResultSetError(err error) bool
	Lookup() LookupRepository
}

type Server struct {
	echo        *echo.Echo
	cfg         config.ServerConfig
	repoManager RepositoryManager
}

func New(cfg config.ServerConfig, repoManager RepositoryManager) *Server {
	e := echo.New()
	e.Server.ReadTimeout = 5 * time.Second
	e.Server.WriteTimeout = 10 * time.Second
	e.Server.IdleTimeout = 120 * time.Second

	return &Server{
		echo:        echo.New(),
		cfg:         cfg,
		repoManager: repoManager,
	}
}

func (s *Server) Run() (err error) {
	s.registerRoutes()

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
