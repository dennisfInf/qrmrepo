package server

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (s *Server) UserAddressMapper() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username := c.Request().Header.Get("x-username")
			lookup, err := s.repoManager.Lookup().GetByUsername(c.Request().Context(), username)
			log.Info().Caller().Msg("received msg")

			if s.repoManager.IsEmptyResultSetError(err) {
				log.Debug().Caller().Err(err).Msg("user accessed endpoint without account")
				return echo.NewHTTPError(http.StatusBadRequest, "you must first sign up")
			} else if err != nil {
				log.Error().Caller().Err(err).Msg("failed to get lookup from database")
				return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}

			c.Set("address", lookup.EnclaveAddress)

			return next(c)
		}
	}
}
