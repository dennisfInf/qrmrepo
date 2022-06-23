package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"

	"log"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

var (
	webAuthn *webauthn.WebAuthn
	session  *webauthn.SessionData
	user     User
)

func BeginRegisterHandler(c echo.Context) error {
	username := c.Request().Header.Get("x-username")
	user = NewUser(username, username)

	log.Println(user)

	log.Printf("received request on: /register/begin with username: %s", username)

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.Parameters = []protocol.CredentialParameter{
			{
				Type:      "public-key",
				Algorithm: -7,
			},
			{
				Type:      "public-key",
				Algorithm: -257,
			},
		}
		credCreationOpts.Attestation = protocol.PreferNoAttestation
		credCreationOpts.CredentialExcludeList = user.CredentialExcludeList()
	}

	var options *protocol.CredentialCreation
	var err error
	options, session, err = webAuthn.BeginRegistration(user, registerOptions)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	log.Println(options)

	return c.JSON(http.StatusOK, options)
}

func FinishRegisterHandler(c echo.Context) error {
	log.Print("received request on: /register/finish")

	credential, err := webAuthn.FinishRegistration(&user, *session, c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	user.AddCredential(*credential)

	log.Println(credential)

	return c.String(http.StatusOK, "register successful")
}

func BeginLoginHandler(c echo.Context) error {
	var options *protocol.CredentialAssertion
	var err error
	options, session, err = webAuthn.BeginLogin(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, options)
}

func FinishLoginHandler(c echo.Context) error {
	credential, err := webAuthn.FinishLogin(user, *session, c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	log.Println(credential)

	return c.String(http.StatusOK, "login successful")
}

func main() {
	var err error

	log.Print("create the webauthn config")
	webAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "enclaive",
		RPID:          "localhost",
		RPOrigin:      "http://localhost",
	})
	if err != nil {
		log.Fatal("failed to create webauthn from config: ", err)
	}

	e := echo.New()
	e.Server.ReadTimeout = 5 * time.Second
	e.Server.WriteTimeout = 10 * time.Second
	e.Server.IdleTimeout = 120 * time.Second

	e.GET("/register/initialize", BeginRegisterHandler)
	e.POST("/register/finalize", FinishRegisterHandler)
	e.GET("/login/initialize", BeginLoginHandler)
	e.POST("/login/finalize", FinishLoginHandler)

	err = e.Start(":80")
	if err == http.ErrServerClosed {
		err = nil
	}
	if err != nil {
		panic(err)
	}
}
