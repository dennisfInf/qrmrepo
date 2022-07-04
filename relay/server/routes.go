package server

import (
	"github.com/enclaive/relay/server/handlers"
)

func (s *Server) registerRoutes() {
	proxy := handlers.ProxyHandler{}

	register := s.echo.Group("/register")
	register.GET("/initialize", proxy.Proxy, s.EnclaveCreator())
	register.POST("/finalize", proxy.Proxy, s.UserAddressMapper())

	login := s.echo.Group("/login", s.UserAddressMapper())
	login.GET("/initialize", proxy.Proxy)
	login.POST("/finalize", proxy.Proxy)

	transaction := s.echo.Group("/transaction", s.UserAddressMapper())
	transaction.GET("/initialize", proxy.Proxy)
	transaction.POST("/finalize", proxy.Proxy)
}
