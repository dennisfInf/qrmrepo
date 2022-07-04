package server

import "github.com/enclaive/backend/server/handlers"

func (s *Server) registerRoutes() error {
	infuraHandler, err := handlers.NewGethHandler(s.cfg.Infura)
	if err != nil {
		return err
	}

	infura := s.echo.Group("/infura")
	infura.GET("/getWalletAddress", infuraHandler.GetWalletAddress())
	infura.POST("/prepareTransaction", infuraHandler.PrepareTransaction())
	infura.POST("/sendTransaction", infuraHandler.SendTransaction())

	return nil
}
