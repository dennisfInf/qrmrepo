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
	"fmt"
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
	Signature [65]byte `json:"signature"`
	ChainID   *big.Int `json:"chain_id"`
	Nonce     uint64   `json:"nonce"`
	GasFeeCap *big.Int `json:"gas_fee_cap"`
	GasTipCap *big.Int `json:"gas_tip_cap"`
	Gas       uint64   `json:"gas"`
	ToAddress string   `json:"to_address"`
	Value     *big.Int `json:"value"`
	Data      []byte   `json:"data"`
}

type prepareTransactionContainer struct {
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

type walletAddress struct {
	Address string `json:"address"`
}

var (
	webAuthn            *webauthn.WebAuthn
	session             *webauthn.SessionData
	user                User
	backendIPAddr       string
	preparedTransaction prepareTransactionContainer
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
		p := point{X: uint64(Cpoint.x), Y: uint64(Cpoint.y)}

		return c.JSON(http.StatusOK, p)
	}
}

func getWalletAddressHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Request().Header.Get("x-username")
		log.Info().Caller().Str("username", username).Msgf("received request on: /getWalletAddress")

		Cpoint := C.host_get_pubkey()
		p := point{X: uint64(Cpoint.x), Y: uint64(Cpoint.y)}

		var address walletAddress
		err := doRequest(http.MethodGet, backendIPAddr+"/getWalletAddress", &p, &address)
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to get wallet address")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		return c.JSON(http.StatusOK, address)
	}
}

// Respond the client with a challenge
func prepareTransactionHandler() echo.HandlerFunc {
	type input struct {
		Amount  uint64 `json:"amount"`
		Address string `json:"address"`
	}

	// Prepare the transaction
	type prepareTransactionIn struct {
		PublicKeyX uint64 `json:"public_key_x"`
		PublicKeyY uint64 `json:"public_key_y"`
		ToAddress  string `json:"to_address"`
		Value      uint64 `json:"value"`
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
		p := point{X: uint64(Cpoint.x), Y: uint64(Cpoint.y)}

		var address walletAddress
		err := doRequest(http.MethodGet, backendIPAddr+"/getWalletAddress", &p, &address)
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to get wallet address")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		prepTranIn := prepareTransactionIn{
			p.X,
			p.Y,
			in.Address,
			in.Amount,
		}

		err = doRequest(http.MethodPost, backendIPAddr+"/prepareTransaction", &prepTranIn, &preparedTransaction)
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to prepare transaction")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		//	options, webSession, err := webAuthn.BeginTransaction(&user, preparedTransaction.Hash)
		//	if err != nil {
		//		log.Error().Caller().Err(err).Msg("failed to prepare webauthn transaction")
		//		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		//	}

		//	session = webSession

		return c.String(http.StatusOK, "prepared successfully")
	}
}

// Receives a request with signed payload incl. challenge
func sendTransactionHandler() echo.HandlerFunc {
	type output struct {
		TransactionHash string `json:"transaction_hash"`
	}

	return func(c echo.Context) error {
		username := c.Request().Header.Get("x-username")
		log.Info().Caller().Str("username", username).Msgf("received request on: /prepareTransaction")

		// Webauthn
		//_, err := webAuthn.FinishTransaction(&user, *session, c.Request())
		//if err != nil {
		//	log.Error().Caller().Err(err).Msg("failed to finish webauthn transaction")
		//	return echo.NewHTTPError(http.StatusBadRequest, "webauthn transaction failed")
		//}

		Csig := C.host_sign_secp256k1((*C.uchar)(C.CBytes(preparedTransaction.Hash[:])), C.uint(len(preparedTransaction.Hash)))
		preparedTransaction.Signature = C.GoBytes(unsafe.Pointer(Csig), 73)

		//TODO fix signature length

		var sigArray [65]byte
		copy(sigArray[:], preparedTransaction.Signature[0:65])

		var out output
		err := doRequest(http.MethodPost, backendIPAddr+"/sendTransaction", &prepareTransactionOut{
			Hash:      preparedTransaction.Hash,
			Signature: sigArray,
			ChainID:   &preparedTransaction.ChainID,
			Nonce:     preparedTransaction.Nonce,
			GasFeeCap: &preparedTransaction.GasFeeCap,
			GasTipCap: &preparedTransaction.GasTipCap,
			Gas:       preparedTransaction.Gas,
			ToAddress: preparedTransaction.ToAddress,
			Value:     &preparedTransaction.Value,
			Data:      preparedTransaction.Data,
		}, &out)
		if err != nil {
			log.Error().Caller().Err(err).Msg("failed to send transaction")
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		log.Info().Caller().Msgf("transaction hash is: %s", out.TransactionHash)

		return c.JSON(http.StatusOK, out)
	}
}

func doRequest(method string, url string, payload any, out any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	log.Info().Caller().Msgf("sending payload %s to %s\n", payloadJSON, url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("couldn't initialize a new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("received error status: %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(out); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
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

	//C.host_gen_secp256k1_keys()

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
	e.POST("/prepareTransaction", prepareTransactionHandler())
	e.POST("/sendTransaction", sendTransactionHandler())

	err = e.Start(":80")
	if err == http.ErrServerClosed {
		err = nil
	}
	if err != nil {
		panic(err)
	}
}
