package server

import (
	"context"
	"fmt"
	"github.com/enclaive/relay/config"
	"github.com/enclaive/relay/models"
	"github.com/labstack/echo/v4"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
	clientset   *kubernetes.Clientset
}

func New(cfg config.ServerConfig, repoManager RepositoryManager) (*Server, error) {
	e := echo.New()
	e.Server.ReadTimeout = 5 * time.Second
	e.Server.WriteTimeout = 10 * time.Second
	e.Server.IdleTimeout = 120 * time.Second

	kubeconf, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconf: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(kubeconf)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &Server{
		echo:        echo.New(),
		cfg:         cfg,
		repoManager: repoManager,
		clientset:   clientset,
	}, nil
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
