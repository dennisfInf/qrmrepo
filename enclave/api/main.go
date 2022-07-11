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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math/big"
	"net/http"
	"os"
	"time"
	"unsafe"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/labstack/echo/v4"
)

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
	X big.Int `json:"public_key_x"`
	Y big.Int `json:"public_key_y"`
}

type walletAddress struct {
	Address string `json:"address"`
}

var (
	webAuthn            *webauthn.WebAuthn
	session             *webauthn.SessionData
	user                User
	backendIPAddr       string
	preparedTransaction prepareTransactionOut
)

func registerInitializeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		if user.startedRegistration() {
			return echo.NewHTTPError(http.StatusBadRequest, "you can't register twice")
		}

		username := c.Request().Header.Get("x-username")
		log.Info().Caller().Str("username", username).Msgf("received request on: /register/initialize")

		user = NewUser(username, username)

		registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
			credCreationOpts.Parameters = []protocol.CredentialParameter{
				{
					Type:      "public-key",
					Algorithm: -7,
				},
			}
			credCreationOpts.Attestation = protocol.PreferNoAttestation
			credCreationOpts.CredentialExcludeList = user.CredentialExcludeList()
		}

		options, sessionData, err := webAuthn.BeginRegistration(user, registerOptions)
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to begin registration")
			return echo.NewHTTPError(http.StatusInternalServerError, "registration failed")
		}

		session = sessionData

		log.Info().Caller().Msgf("sending response: %v", options)

		return c.JSON(http.StatusOK, options)
	}
}

func registerFinalizeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		if user.finishedRegistration() {
			return echo.NewHTTPError(http.StatusBadRequest, "you can't register twice")
		}

		username := c.Request().Header.Get("x-username")
		log.Info().Caller().Str("username", username).Msgf("received request on: /register/finalize")

		credential, err := webAuthn.FinishRegistration(&user, *session, c.Request())
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to finish registration")
			return echo.NewHTTPError(http.StatusBadRequest, "registration failed")
		}

		log.Info().Caller().Msgf("adding credential: %v", credential)
		user.AddCredential(*credential)

		log.Info().Caller().Msg("saving user")
		err = SaveUser(user)
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to save user")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		return c.String(http.StatusOK, "register successful")
	}
}

func loginInitializeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Request().Header.Get("x-username")
		log.Info().Caller().Str("username", username).Msgf("received request on: /login/initialize")

		options, webSession, err := webAuthn.BeginLogin(&user)
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to begin login")
			return echo.NewHTTPError(http.StatusBadRequest, "login failed")
		}

		session = webSession

		log.Info().Caller().Msgf("sending response: %v", options)

		return c.JSON(http.StatusOK, options)
	}
}

func loginFinalizeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Request().Header.Get("x-username")
		log.Info().Caller().Str("username", username).Msgf("received request on: /login/finalize")

		_, err := webAuthn.FinishLogin(&user, *session, c.Request())
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to finish login")
			return echo.NewHTTPError(http.StatusBadRequest, "login failed")
		}

		return c.String(http.StatusOK, "login successful")
	}
}

func getPublicKeyHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Request().Header.Get("x-username")
		log.Info().Caller().Str("username", username).Msgf("received request on: /getPublicKey")

		Cpoint := C.host_get_pubkey()
		p := point{
			X: Cpoint.x,
			Y: Cpoint.y,
		}

		return c.JSON(http.StatusOK, p)
	}
}

func getWalletAddressHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Request().Header.Get("x-username")
		log.Info().Caller().Str("username", username).Msgf("received request on: /getWalletAddress")

		Cpoint := C.host_get_pubkey()
		p := point{
			X: Cpoint.x,
			Y: Cpoint.y,
		}

		pointJSON, err := json.Marshal(p)
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to marshal point")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		log.Info().Caller().Msgf("requesting address for point: %s", pointJSON)
		req, err := http.NewRequest("GET", backendIPAddr+"/getWalletAddress", bytes.NewBuffer(pointJSON))
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to make new request")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error().Caller().Err(err).Msg("request to backend failed")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
		defer res.Body.Close()

		var address walletAddress
		if err := json.NewDecoder(res.Body).Decode(&address); err != nil {
			log.Error().Caller().Err(err).Msg("failed to decode backend response")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		return c.JSON(http.StatusOK, address)
	}
}

// Respond the client with a challenge
func prepareTransactionHandler() echo.HandlerFunc {
	type input struct {
		Amount  uint   `json:"amount"`
		Address []byte `json:"address"`
	}

	// Prepare the transaction
	type prepareTransactionIn struct {
		PublicKeyX big.Int `json:"public_key_x"`
		PublicKeyY big.Int `json:"public_key_y"`
		ToAddress  string  `json:"to_address"`
		Value      uint    `json:"value"`
	}

	return func(c echo.Context) error {
		username := c.Request().Header.Get("x-username")
		log.Info().Caller().Str("username", username).Msgf("received request on: /prepareTransaction")

		var in input
		if err := c.Bind(&in); err != nil {
			log.Error().Caller().Err(err).Msg("failed to read input")
			return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
		}

		// Get the public key from the enclave
		Cpoint := C.host_get_pubkey()
		p := point{
			X: Cpoint.x,
			Y: Cpoint.y,
		}

		pointJSON, err := json.Marshal(p)
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to marshal point")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		log.Info().Caller().Msgf("requesting address for point: %s", pointJSON)
		req, err := http.NewRequest("GET", backendIPAddr+"/getWalletAddress", bytes.NewBuffer(pointJSON))
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to make new request")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error().Caller().Err(err).Msg("request to backend failed")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
		defer res.Body.Close()

		var address walletAddress
		if err := json.NewDecoder(res.Body).Decode(&address); err != nil {
			log.Error().Caller().Err(err).Msg("failed to decode backend response")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		prepTranIn := prepareTransactionIn{
			p.X,
			p.Y,
			string(in.Address),
			in.Amount,
		}

		prepTranInJSON, err := json.Marshal(prepTranIn)
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to marshal prepareTransactionIn")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		reqPrepTran, err := http.NewRequest("POST", backendIPAddr+"/prepareTransaction", bytes.NewBuffer(prepTranInJSON))
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to create prepareTransaction request")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
		reqPrepTran.Header.Set("Content-Type", "application/json")

		res, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Error().Caller().Err(err).Msg("request to backend failed")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
		defer res.Body.Close()

		var preparedTransaction prepareTransactionOut
		if err := json.NewDecoder(res.Body).Decode(&preparedTransaction); err != nil {
			log.Error().Caller().Err(err).Msg("failed to decode backend response")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		options, webSession, err := webAuthn.BeginTransaction(&user, preparedTransaction.Hash)
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to prepare webauthn transaction")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		session = webSession

		return c.JSON(http.StatusOK, options)
	}
}

// Receives a request with signed payload incl. challenge
// {address: "...", value: xx, nonce: "..."}
// +
// Signature
func sendTransactionHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Request().Header.Get("x-username")
		log.Info().Caller().Str("username", username).Msgf("received request on: /prepareTransaction")

		// Webauthn
		_, err := webAuthn.FinishTransaction(&user, *session, c.Request())
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to finish webauthn transaction")
			return echo.NewHTTPError(http.StatusBadRequest, "webauthn transaction failed")
		}

		Csig := C.host_sign_secp256k1((*C.uchar)(C.CBytes(preparedTransaction.Hash[:])), C.uint(len(preparedTransaction.Hash)))
		preparedTransaction.Signature = C.GoBytes(unsafe.Pointer(Csig), 73)

		preparedTransactionJSON, err := json.Marshal(preparedTransaction)
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to marshal prepareTransactionIn")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		reqPrepTran, err := http.NewRequest("POST", backendIPAddr+"/sendTransaction", bytes.NewBuffer(preparedTransactionJSON))
		if err != nil {
			log.Printf("Couldn't initialize a new request, %s\n", err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}
		reqPrepTran.Header.Set("Content-Type", "application/preparedTransactionJSON")

		res, err := http.DefaultClient.Do(reqPrepTran)
		if err != nil {
			log.Error().Caller().Err(err).Msg("request to backend failed")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
		defer res.Body.Close()

		return c.String(http.StatusOK, "transaction successful")
	}
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	var err error
	user, err = LoadUser()
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("failed to load user")
	}

	backendIPAddr = "http://" + os.Getenv("BACKEND_IP") + "/infura"

	C.host_gen_secp256k1_keys()

	//log.Print("==== Testing Sign ====")
	//C.test_sign_secp256k1()
	//log.Print("==== End Testing ====")

	log.Info().Caller().Msg("creating the webauthn config")
	webAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "enclaive",
		RPID:          "elonmaskwallet.com",
		RPOrigin:      "https://elonmaskwallet.com",
	})
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("failed to create webauthn config")
	}

	e := echo.New()
	e.Server.ReadTimeout = 5 * time.Second
	e.Server.WriteTimeout = 10 * time.Second
	e.Server.IdleTimeout = 120 * time.Second

	e.GET("/register/initialize", registerInitializeHandler())
	e.POST("/register/finalize", registerFinalizeHandler())
	e.GET("/login/initialize", loginInitializeHandler())
	e.POST("/login/finalize", loginFinalizeHandler())
	e.GET("/getPublicKey", getPublicKeyHandler())
	e.GET("/getWalletAddress", getWalletAddressHandler())
	e.GET("/prepareTransaction", prepareTransactionHandler())
	e.GET("/sendTransaction", sendTransactionHandler())

	err = e.Start(":80")
	if err == http.ErrServerClosed {
		err = nil
	}
	if err != nil {
		panic(err)
	}
}
