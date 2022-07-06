package main

/*
#cgo CFLAGS: -I../host/
#cgo LDFLAGS: -L../build/host/ -L/opt/openenclave/lib/openenclave/host -Wl,-z,noexecstack -lhost -loehost -ldl -lpthread -lssl -lcrypto -lstdc++

#include "host.h"
*/
import "C"

import (
	"encoding/json"
	"net/http"

	"log"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-gonic/gin"
)

var webAuthn *webauthn.WebAuthn
var session *webauthn.SessionData
var user webAuthnUser

type webAuthnUser struct {
	id         []byte
	name       string
	User       string
	credential *webauthn.Credential
}

func (user *webAuthnUser) WebAuthnID() []byte {
	return user.id
}
func (user *webAuthnUser) WebAuthnName() string {
	return user.name
}
func (user *webAuthnUser) WebAuthnDisplayName() string {
	return user.name
}
func (user *webAuthnUser) WebAuthnIcon() string {
	return ""
}
func (user *webAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return []webauthn.Credential{*user.credential}
}

func BeginRegisterHandler(c *gin.Context) {
	//username := c.GetHeader("X-Authenticated-User")
	username := c.Param("user")
	log.Printf("received request on: /register/begin with username: %s", username)

	user.name = username
	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		// TODO
		credCreationOpts.Attestation = protocol.PreferIndirectAttestation
	}

	options, sessionData, err := webAuthn.BeginRegistration(&webAuthnUser{[]byte(""), username, username, nil}, registerOptions)
	session = sessionData

	optionsJSON, _ := json.Marshal(options)
	log.Printf("Relying Party response to client: %s\n", optionsJSON)

	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONP(http.StatusOK, options)
}

func FinishRegisterHandler(c *gin.Context) {
	log.Print("received request on: /register/finish")

	credential, err := webAuthn.FinishRegistration(&user, *session, c.Request)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	user.credential = credential

	credentialJSON, _ := json.MarshalIndent(credential, "", " ")

	log.Printf("Credential: %+v\n", credentialJSON)
	c.Data(http.StatusOK, "text/html", []byte(""))

}

func BeginLoginHandler(c *gin.Context) {
	log.Print("received request on: /login/begin")
	log.Print(user)

	options, webSession, err := webAuthn.BeginLogin(&user)
	session = webSession

	log.Print(session)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, options)
}

func FinishLogin(c *gin.Context) {
	log.Print("received request on: /login/finish")

	_, err := webAuthn.FinishLogin(&user, *session, c.Request)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Login Success")
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

	r := gin.Default()
	r.Use(gin.Logger())
	r.GET("/register/begin/:user", BeginRegisterHandler)
	r.POST("/register/finish/:user", FinishRegisterHandler)
	r.GET("/login/begin/:user", BeginLoginHandler)
	r.POST("/login/finish/:user", FinishLogin)
	r.Static("/static", "../")
	r.Run(":2533")
}
