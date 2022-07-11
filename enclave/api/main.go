package main

/*
#cgo CFLAGS: -I../host/
#cgo LDFLAGS: -L../build/host/ -L/opt/openenclave/lib/openenclave/host -Wl,-z,noexecstack -lhost -loehost -ldl -lpthread -lssl -lcrypto -lstdc++

#include "host.h"
#include <stdlib.h>
*/
import "C"

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"time"
	"unsafe"

	"log"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/labstack/echo/v4"
)

var webAuthn *webauthn.WebAuthn
var session *webauthn.SessionData
var user webAuthnUser
var backendIPAddr string
var preparedTransaction prepareTransactionOut

type prepareTransactionOut struct {
	Hash      [32]byte `json:"hash"`
	Signature []byte   `json:"signature"`
	ChainID   big.Int  `json:"chain_id"`
	Nonce     uint64   `json:"nonce"`
	GasFeeCap big.Int  `json:"gas_fee_cap"`
	GasTipCap big.Int  `json:"gas_tip_cap"`
	Gas       uint64   `json:"gas"`
	ToAddress string   `json:"to_address"`
	Value     big.Int  `json:"value"`
	Data      []byte   `json:"data"`
}
type point struct {
	X uint64 `json:"public_key_x"`
	Y uint64 `json:"public_key_y"`
}

type challenge struct {
	Challenge []byte `json:"challenge"`
}

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

func BeginRegisterHandler(c echo.Context) error {
	username := c.Request().Header.Get("x-username")
	log.Printf("received request on: /register/initialize with username: %s", username)

	user.name = username
	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.Attestation = protocol.PreferNoAttestation
	}

	options, sessionData, err := webAuthn.BeginRegistration(&webAuthnUser{[]byte(""), username, username, nil}, registerOptions)
	session = sessionData

	optionsJSON, _ := json.Marshal(options)
	log.Printf("Relying Party response to client: %s\n", optionsJSON)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, options)
}

func FinishRegisterHandler(c echo.Context) error {
	log.Print("received request on: /register/finalize")

	credential, err := webAuthn.FinishRegistration(&user, *session, c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	user.credential = credential

	log.Println(credential)

	return c.String(http.StatusOK, "register successful")
}

func BeginLoginHandler(c echo.Context) error {
	log.Print("received request on: /login/initialize")
	log.Print(user)

	options, webSession, err := webAuthn.BeginLogin(&user)
	session = webSession

	log.Print(session)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, options)
}

func FinishLoginHandler(c echo.Context) error {
	log.Print("received request on: /login/finalize")

	_, err := webAuthn.FinishLogin(&user, *session, c.Request())
	if err != nil {
		log.Print(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "Login Success")
}

func getPublicKey(c echo.Context) error {
	log.Print("received request on: /getPublicKey")
	Cpoint := C.host_get_pubkey()
	p := point{X: uint64(Cpoint.x), Y: uint64(Cpoint.y)}
	return c.JSON(http.StatusOK, p)
}

func getWalletAddress(c echo.Context) error {
	log.Print("received request on: /getWalletAddress")
	Cpoint := C.host_get_pubkey()
	p := point{X: uint64(Cpoint.x), Y: uint64(Cpoint.y)}
	js, err := json.Marshal(p)
	log.Printf("JSON: %s", js)
	if err != nil {
		log.Printf("Couldn't marshal the struct, %s\n", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	req, err := http.NewRequest("GET", backendIPAddr+"/getWalletAddress", bytes.NewBuffer(js))
	if err != nil {
		log.Printf("Couldn't initialize a new request, %s\n", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error occured during request, %s\n", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return c.JSONBlob(http.StatusOK, body)
}

// Response the client with a challenge
func beginTransaction(c echo.Context) error {
	log.Print("received request on /beginTransaction")

	type input struct {
		Amount  uint   `json:"amount"`
		Address []byte `json:"address"`
	}

	var req input

	if err := c.Bind(&req); err != nil {
		log.Printf("Couldn't bind the json, %s", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Get the publix key from the enclave
	Cpoint := C.host_get_pubkey()
	p := point{X: uint64(Cpoint.x), Y: uint64(Cpoint.y)}
	js, err := json.Marshal(p)
	if err != nil {
		log.Printf("Couldn't marshal the struct, %s\n", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Get the wallet address from the backend
	reqWallet, err := http.NewRequest("GET", backendIPAddr+"/getWalletAddress", bytes.NewBuffer(js))
	if err != nil {
		log.Printf("Couldn't initialize a new request, %s\n", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	reqWallet.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqWallet)
	if err != nil {
		log.Printf("Error occured during request, %s\n", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	type output struct {
		Address string `json:"address"`
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error occured during parsing response, %s\n", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var respWallet output
	json.Unmarshal(body, respWallet)

	// Prepare the transaction
	type prepareTransactionIn struct {
		PublicKeyX uint64 `json:"public_key_x"`
		PublicKeyY uint64 `json:"public_key_y"`
		ToAddress  string `json:"to_address"`
		Value      uint   `json:"value"`
	}

	prepTranIn := prepareTransactionIn{p.X, p.Y, string(req.Address), req.Amount}
	prepTranInjs, err := json.Marshal(prepTranIn)

	if err != nil {
		log.Printf("Couldn't marshal the struct, %s\n", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	reqPrepTran, err := http.NewRequest("POST", backendIPAddr+"/prepareTransaction", bytes.NewBuffer(prepTranInjs))
	if err != nil {
		log.Printf("Couldn't initialize a new request, %s\n", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	reqWallet.Header.Set("Content-Type", "application/json")

	client01 := &http.Client{}
	respPrepTran, err := client01.Do(reqPrepTran)
	if err != nil {
		log.Printf("Error occured during request, %s\n", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer respPrepTran.Body.Close()

	bodyRespPrepTran, err := ioutil.ReadAll(respPrepTran.Body)

	if err != nil {
		log.Printf("Error occured during parsing response, %s\n", err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	json.Unmarshal(bodyRespPrepTran, preparedTransaction)

	// Webauthn
	options, webSession, err := webAuthn.BeginTransaction(&user, preparedTransaction.Hash)
	session = webSession

	log.Print(session)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, options)
}

// Receives a request with signed payload incl. challenge
// {address: "...", value: xx, nonce: "..."}
// +
// Signature
func finishTransaction(c echo.Context) error {
	log.Print("received request on /finishTransaction")

	/*_, err := webAuthn.FinishLogin(&user, *session, c.Request)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}*/

	Csig := C.host_sign_secp256k1((*C.uchar)(C.CBytes(preparedTransaction.Hash[:])), C.uint(len(preparedTransaction.Hash)))

	preparedTransaction.Signature = C.GoBytes(unsafe.Pointer(Csig), 73)

	json, _ := json.Marshal(preparedTransaction)

	reqPrepTran, err := http.NewRequest("POST", backendIPAddr+"/sendTransaction", bytes.NewBuffer(json))
	if err != nil {
		log.Printf("Couldn't initialize a new request, %s\n", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	reqPrepTran.Header.Set("Content-Type", "application/json")

	client01 := &http.Client{}
	respPrepTran, err := client01.Do(reqPrepTran)
	if err != nil {
		log.Printf("Error occured during request, %s\n", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer respPrepTran.Body.Close()

	return c.String(http.StatusOK, "Login Success")
}

func main() {
	var err error

	backendIPAddr = "http://" + os.Getenv("BACKEND_IP") + "/infura"

	C.host_gen_secp256k1_keys()

	//log.Print("==== Testing Sign ====")
	//C.test_sign_secp256k1()
	//log.Print("==== End Testing ====")

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
	e.GET("/getPublicKey", getPublicKey)
	e.GET("/getWalletAddress", getWalletAddress)
	e.GET("/prepareTransaction", beginTransaction)
	e.GET("/sendTransaction", finishTransaction)

	err = e.Start(":80")
	if err == http.ErrServerClosed {
		err = nil
	}
	if err != nil {
		panic(err)
	}
}
