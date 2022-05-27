package server

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type JWTCustomClaims struct {
	EnclaveAddress string `json:"enclave_address"`
	jwt.StandardClaims
}

func Session(repoMgr RepositoryManager, signingSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("x-relay-token")

			if tokenString == "" { //Create new Session
				var err error

				tokenString, err = createSession(c, repoMgr, signingSecret)
				if err != nil {
					log.Error().Caller().Err(err).Msg("failed to create session")
					return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				}
			}

			//Validate session
			token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(signingSecret), nil
			})
			if err != nil {
				log.Info().Caller().Err(err).Msg("jwt could not be parsed")
				return c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}

			claims, ok := token.Claims.(*JWTCustomClaims)
			if !ok || !token.Valid {
				return c.String(http.StatusUnauthorized, "jwt is invalid")
			}

			c.Set("address", claims.EnclaveAddress)

			return next(c)
		}
	}
}

func createSession(c echo.Context, repoMgr RepositoryManager, signingSecret string) (string, error) {
	username := c.Request().Header.Get("x-username")
	lookup, err := repoMgr.Lookup().GetByUsername(c.Request().Context(), username)

	if err != nil {
		if repoMgr.IsEmptyResultSetError(err) {
			//TODO spawn enclave
			return "", fmt.Errorf("failed to spawn enclave: %w", err)
		} else {
			return "", fmt.Errorf("failed to get lookup from db: %w", err)
		}
	}

	token, err := createSessionToken(lookup.EnclaveAddress, signingSecret)
	if err != nil {
		return "", fmt.Errorf("failed to create session token: %w", err)
	}

	// Set header so that the client can use the token for the coming requests
	c.Response().Header().Set("x-set-relay-token", token)

	return token, nil
}

func createSessionToken(enclaveAddress, signingSecret string) (string, error) {
	currentTime := time.Now().UTC()

	claims := &JWTCustomClaims{
		EnclaveAddress: enclaveAddress,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: currentTime.Add(time.Hour * 24).Unix(),
			IssuedAt:  currentTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signingSecret))
}
