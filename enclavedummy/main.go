package main

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
	id   []byte
	name string
	User string
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
	return []webauthn.Credential{}
}

func BeginRegisterHandler(c *gin.Context) {
	//username := c.GetHeader("X-Authenticated-User")
	username := c.GetHeader("x-username")
	log.Printf("received request on: /register/begin with username: %s", username)

	user.name = username
	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		// TODO
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
	}

	options, sessionData, err := webAuthn.BeginRegistration(&webAuthnUser{[]byte("1"), username, ""}, registerOptions)
	session = sessionData

	optionsJSON, _ := json.Marshal(options)
	log.Printf("Relying Party response to client: %s\n", optionsJSON)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONP(http.StatusOK, options)
}

func FinishRegisterHandler(c *gin.Context) {
	log.Print("received request on: /register/finish")

	credential, err := webAuthn.FinishRegistration(&user, *session, c.Request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	credentialJSON, _ := json.MarshalIndent(credential, "", " ")

	log.Printf("Credential: %s\n", credentialJSON)
	c.Data(http.StatusOK, "text/html", []byte(""))

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
	r.GET("/register/initialize", BeginRegisterHandler)
	r.POST("/register/finalize", FinishRegisterHandler)
	r.Static("/static", "../")
	r.Run(":8080")
}
