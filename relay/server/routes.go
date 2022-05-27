package server

import (
	"github.com/enclaive/relay/server/handlers"
)

func (s *Server) registerRoutes() {
	proxy := handlers.ProxyHandler{}

	register := s.echo.Group("/register", UserAddressMapper(s.repoManager))
	register.POST("/initialize", proxy.Proxy)
	register.POST("/finalize", proxy.Proxy)

	login := s.echo.Group("/login", UserAddressMapper(s.repoManager))
	login.POST("/initialize", proxy.Proxy)
	login.POST("/finalize", proxy.Proxy)
}
