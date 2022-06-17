package server

import (
	"github.com/enclaive/relay/server/handlers"
)

func (s *Server) registerRoutes() {
	proxy := handlers.ProxyHandler{}

	register := s.echo.Group("/register", s.UserAddressMapper())
	register.GET("/initialize", proxy.Proxy)
	register.POST("/finalize", proxy.Proxy)

	login := s.echo.Group("/login", s.UserAddressMapper())
	login.GET("/initialize", proxy.Proxy)
	login.POST("/finalize", proxy.Proxy)
}
