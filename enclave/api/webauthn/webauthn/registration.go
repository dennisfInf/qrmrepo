package webauthn

/*
#cgo CFLAGS: -I/server/host/
#cgo LDFLAGS: -L/server/build/host/ -L/opt/openenclave/lib/openenclave/host -Wl,-z,noexecstack -lhost -loehost -ldl -lpthread -lssl -lcrypto -lstdc++

#include "host.h"
*/
import "C"

import (
	"bytes"
	"encoding/base64"
	"log"
	"math/big"
	"net/http"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/protocol/webauthncose"
)

// BEGIN REGISTRATION
// These objects help us creat the CredentialCreationOptions
// that will be passed to the authenticator via the user client

type RegistrationOption func(*protocol.PublicKeyCredentialCreationOptions)

// Generate a new set of registration data to be sent to the client and authenticator.
func (webauthn *WebAuthn) BeginRegistration(user User, opts ...RegistrationOption) (*protocol.CredentialCreation, *SessionData, error) {
	challenge, err := protocol.CreateChallenge()
	if err != nil {
		return nil, nil, err
	}

	webAuthnUser := protocol.UserEntity{
		ID:          user.WebAuthnID(),
		DisplayName: user.WebAuthnDisplayName(),
		CredentialEntity: protocol.CredentialEntity{
			Name: user.WebAuthnName(),
			Icon: user.WebAuthnIcon(),
		},
	}

	relyingParty := protocol.RelyingPartyEntity{
		ID: webauthn.Config.RPID,
		CredentialEntity: protocol.CredentialEntity{
			Name: webauthn.Config.RPDisplayName,
			Icon: webauthn.Config.RPIcon,
		},
	}

	credentialParams := defaultRegistrationCredentialParameters()

	creationOptions := protocol.PublicKeyCredentialCreationOptions{
		Challenge:              challenge,
		RelyingParty:           relyingParty,
		User:                   webAuthnUser,
		Parameters:             credentialParams,
		AuthenticatorSelection: webauthn.Config.AuthenticatorSelection,
		Timeout:                webauthn.Config.Timeout,
		Attestation:            webauthn.Config.AttestationPreference,
	}

	for _, setter := range opts {
		setter(&creationOptions)
	}

	response := protocol.CredentialCreation{Response: creationOptions}
	newSessionData := SessionData{
		Challenge:        base64.RawURLEncoding.EncodeToString(challenge),
		UserID:           user.WebAuthnID(),
		UserVerification: creationOptions.AuthenticatorSelection.UserVerification,
	}

	if err != nil {
		return nil, nil, protocol.ErrParsingData.WithDetails("Error packing session data")
	}

	return &response, &newSessionData, nil
}

// Provide non-default parameters regarding the authenticator to select.
func WithAuthenticatorSelection(authenticatorSelection protocol.AuthenticatorSelection) RegistrationOption {
	return func(cco *protocol.PublicKeyCredentialCreationOptions) {
		cco.AuthenticatorSelection = authenticatorSelection
	}
}

// Provide non-default parameters regarding credentials to exclude from retrieval.
func WithExclusions(excludeList []protocol.CredentialDescriptor) RegistrationOption {
	return func(cco *protocol.PublicKeyCredentialCreationOptions) {
		cco.CredentialExcludeList = excludeList
	}
}

// Provide non-default parameters regarding whether the authenticator should attest to the credential.
func WithConveyancePreference(preference protocol.ConveyancePreference) RegistrationOption {
	return func(cco *protocol.PublicKeyCredentialCreationOptions) {
		cco.Attestation = preference
	}
}

// Provide extension parameter to registration options
func WithExtensions(extension protocol.AuthenticationExtensions) RegistrationOption {
	return func(cco *protocol.PublicKeyCredentialCreationOptions) {
		cco.Extensions = extension
	}
}

// WithResidentKeyRequirement sets both the resident key and require resident key protocol options. This could conflict
// with webauthn.WithAuthenticatorSelection if it doesn't come after it.
func WithResidentKeyRequirement(requirement protocol.ResidentKeyRequirement) RegistrationOption {
	return func(cco *protocol.PublicKeyCredentialCreationOptions) {
		cco.AuthenticatorSelection.ResidentKey = requirement
		switch requirement {
		case protocol.ResidentKeyRequirementRequired:
			cco.AuthenticatorSelection.RequireResidentKey = protocol.ResidentKeyRequired()
		default:
			cco.AuthenticatorSelection.RequireResidentKey = protocol.ResidentKeyUnrequired()
		}
	}
}

// Take the response from the authenticator and client and verify the credential against the user's credentials and
// session data.
func (webauthn *WebAuthn) FinishRegistration(user User, session SessionData, response *http.Request) (*Credential, error) {
	parsedResponse, err := protocol.ParseCredentialCreationResponse(response)
	if err != nil {
		return nil, err
	}

	return webauthn.CreateCredential(user, session, parsedResponse)
}

// CreateCredential verifies a parsed response against the user's credentials and session data.
func (webauthn *WebAuthn) CreateCredential(user User, session SessionData, parsedResponse *protocol.ParsedCredentialCreationData) (*Credential, error) {
	if !bytes.Equal(user.WebAuthnID(), session.UserID) {
		return nil, protocol.ErrBadRequest.WithDetails("ID mismatch for User and Session")
	}

	shouldVerifyUser := session.UserVerification == protocol.VerificationRequired

	invalidErr := parsedResponse.Verify(session.Challenge, shouldVerifyUser, webauthn.Config.RPID, webauthn.Config.RPOrigin)
	if invalidErr != nil {
		return nil, invalidErr
	}

	// Pass the public key to the enclave
	pkey, err := webauthncose.ParsePublicKey(parsedResponse.Response.AttestationObject.AuthData.AttData.CredentialPublicKey)
	if err != nil {
		log.Printf("Error occured while parsing the public key")
	} else {
		log.Println(pkey)
	}

	switch pkey.(type) {
	case webauthncose.EC2PublicKeyData:
		e := pkey.(webauthncose.EC2PublicKeyData)
		if webauthncose.COSEAlgorithmIdentifier(e.Algorithm) == webauthncose.AlgES256 {
			log.Printf("AlgES256 curve")

			//log.Printf("Elliptic curve public key")
			//log.Printf("Public Key Length: %d\n", len(parsedResponse.Response.AttestationObject.AuthData.AttData.CredentialPublicKey))
			//str, _ := base64.StdEncoding.DecodeString(base64.StdEncoding.EncodeToString(parsedResponse.Response.AttestationObject.AuthData.AttData.CredentialPublicKey))
			x := big.NewInt(0).SetBytes(e.XCoord)
			y := big.NewInt(0).SetBytes(e.YCoord)

			log.Printf("XCoord: %d, len: %d\nYCoord: %d, len: %d\n", x, x.BitLen(), y, y.BitLen())
			log.Printf("CALL THE ENCLAVE\n")
			C.host_store_ecc_pk((*C.uchar)(C.CBytes(e.XCoord)), (*C.uchar)(C.CBytes(e.YCoord)))

			//log.Printf("Public Key: %s\n", str)
			//arr, _ := asn1.Unmarshal(parsedResponse.Response.AttestationObject.AuthData.AttData.CredentialPublicKey, nil)
			//log.Printf("UNMARSHALLED: %s", arr)
		}
		break
	default:
		log.Printf("Unsupported public key type by enclave")
	}

	return MakeNewCredential(parsedResponse)
}

func defaultRegistrationCredentialParameters() []protocol.CredentialParameter {
	return []protocol.CredentialParameter{
		{
			Type:      protocol.PublicKeyCredentialType,
			Algorithm: webauthncose.AlgES256,
		},
		/*{
			Type:      protocol.PublicKeyCredentialType,
			Algorithm: webauthncose.AlgES384,
		},
		{
			Type:      protocol.PublicKeyCredentialType,
			Algorithm: webauthncose.AlgES512,
		},
		{
			Type:      protocol.PublicKeyCredentialType,
			Algorithm: webauthncose.AlgRS256,
		},
		{
			Type:      protocol.PublicKeyCredentialType,
			Algorithm: webauthncose.AlgRS384,
		},
		{
			Type:      protocol.PublicKeyCredentialType,
			Algorithm: webauthncose.AlgRS512,
		},
		{
			Type:      protocol.PublicKeyCredentialType,
			Algorithm: webauthncose.AlgPS256,
		},
		{
			Type:      protocol.PublicKeyCredentialType,
			Algorithm: webauthncose.AlgPS384,
		},
		{
			Type:      protocol.PublicKeyCredentialType,
			Algorithm: webauthncose.AlgPS512,
		},
		{
			Type:      protocol.PublicKeyCredentialType,
			Algorithm: webauthncose.AlgEdDSA,
		},*/
	}
}
