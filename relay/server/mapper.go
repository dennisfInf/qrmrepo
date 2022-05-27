package server

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

func UserAddressMapper(repoManager RepositoryManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username := c.Request().Header.Get("x-username")
			lookup, err := repoManager.Lookup().GetByUsername(c.Request().Context(), username)

			if err != nil {
				if repoManager.IsEmptyResultSetError(err) {
					//TODO spawn enclave
					log.Error().Caller().Err(err).Msg("failed to spawn enclave")
				} else {
					log.Error().Caller().Err(err).Msg("failed to get lookup from database")
				}
				return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}

			c.Set("address", lookup.EnclaveAddress)

			return next(c)
		}
	}
}
