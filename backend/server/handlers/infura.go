package handlers

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/enclaive/backend/config"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"math/big"
	"net/http"
)

type InfuraHandler struct {
	client *ethclient.Client
	cfg    config.InfuraConfig
}

func NewGethHandler(cfg config.InfuraConfig) (*InfuraHandler, error) {
	client, err := ethclient.Dial(cfg.InfuraAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to infura: %w", err)
	}

	return &InfuraHandler{
		client: client,
	}, nil
}

func (w *InfuraHandler) GetWalletAddress() echo.HandlerFunc {
	type input struct {
		PublicKeyX big.Int `json:"public_key_x"`
		PublicKeyY big.Int `json:"public_key_y"`
	}

	type output struct {
		Address string `json:"address"`
	}

	return func(c echo.Context) error {
		var in input
		if err := c.Bind(&in); err != nil {
			return c.String(http.StatusBadRequest, "invalid json")
		}

		pubkey := ecdsa.PublicKey{
			Curve: secp256k1.S256(),
			X:     &in.PublicKeyX,
			Y:     &in.PublicKeyY,
		}

		address := crypto.PubkeyToAddress(pubkey)

		return c.JSON(http.StatusOK, output{
			Address: address.String(),
		})
	}
}

func (w *InfuraHandler) PrepareTransaction() echo.HandlerFunc {
	type input struct {
		PublicKeyX big.Int `json:"public_key_x"`
		PublicKeyY big.Int `json:"public_key_y"`
		ToAddress  string  `json:"to_address"`
		Value      big.Int `json:"value"`
	}

	type output struct {
		Hash      [32]byte `json:"hash"`
		ChainID   big.Int  `json:"chain_id"`
		Nonce     uint64   `json:"nonce"`
		GasFeeCap big.Int  `json:"gas_fee_cap"`
		GasTipCap big.Int  `json:"gas_tip_cap"`
		Gas       uint64   `json:"gas"`
		ToAddress string   `json:"to_address"`
		Value     big.Int  `json:"value"`
		Data      []byte   `json:"data"`
	}

	return func(c echo.Context) error {
		var in input
		if err := c.Bind(&in); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid json")
		}

		pubkey := ecdsa.PublicKey{
			Curve: secp256k1.S256(),
			X:     &in.PublicKeyX,
			Y:     &in.PublicKeyY,
		}

		fromAddress := crypto.PubkeyToAddress(pubkey)

		nonce, err := w.client.PendingNonceAt(c.Request().Context(), fromAddress)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		toAddress := common.HexToAddress(in.ToAddress)
		var data []byte

		estimatedGas, err := w.client.EstimateGas(c.Request().Context(), ethereum.CallMsg{
			To:   &toAddress,
			Data: []byte{0},
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		gasLimit := uint64(float64(estimatedGas) * 1.30) // in units
		tipCap := big.NewInt(20000000000)                // maxPriorityFeePerGas = 20 Gwei
		feeCap := big.NewInt(200000000000)               // maxFeePerGas = 200 Gwei

		chainID, err := w.client.NetworkID(c.Request().Context())
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			GasFeeCap: feeCap,
			GasTipCap: tipCap,
			Gas:       gasLimit,
			To:        &toAddress,
			Value:     &in.Value,
			Data:      data,
		})

		h := types.LatestSignerForChainID(chainID).Hash(tx)

		return c.JSON(http.StatusOK, output{
			Hash:      h,
			ChainID:   *chainID,
			Nonce:     nonce,
			GasFeeCap: *feeCap,
			GasTipCap: *tipCap,
			Gas:       gasLimit,
			ToAddress: in.ToAddress,
			Value:     in.Value,
			Data:      data,
		})
	}
}

func (w *InfuraHandler) SendTransaction() echo.HandlerFunc {
	type input struct {
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

	return func(c echo.Context) error {
		var in input
		if err := c.Bind(&in); err != nil {
			return c.String(http.StatusBadRequest, "invalid json")
		}

		toAddress := common.HexToAddress(in.ToAddress)

		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   &in.ChainID,
			Nonce:     in.Nonce,
			GasFeeCap: &in.GasFeeCap,
			GasTipCap: &in.GasTipCap,
			Gas:       in.Gas,
			To:        &toAddress,
			Value:     &in.Value,
			Data:      in.Data,
		})

		signedTx, err := tx.WithSignature(types.LatestSignerForChainID(&in.ChainID), in.Signature)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return w.client.SendTransaction(c.Request().Context(), signedTx)
	}
}
