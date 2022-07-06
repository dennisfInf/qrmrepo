package protocol

/*
#cgo CFLAGS: -I/server/host/
#cgo LDFLAGS: -L/server/build/host/ -L/opt/openenclave/lib/openenclave/host -Wl,-z,noexecstack -lhost -loehost -ldl -lpthread -lssl -lcrypto -lstdc++

#include "host.h"
#include <stdlib.h>
*/
import "C"

import (
	"encoding/base64"
	"unsafe"
)

// ChallengeLength - Length of bytes to generate for a challenge
const ChallengeLength = 32

// Challenge that should be signed and returned by the authenticator
type Challenge URLEncodedBase64

// Create a new challenge to be sent to the authenticator. The spec recommends using
// at least 16 bytes with 100 bits of entropy. We use 32 bytes.
func CreateChallenge() (Challenge, error) {
	//challenge := make([]byte, ChallengeLength)

	// Receive a nonce from the enclave
	Cchallenge := C.host_create_nonce(ChallengeLength)
	challenge := C.GoBytes(unsafe.Pointer(Cchallenge), C.int(ChallengeLength))
	C.free(unsafe.Pointer(Cchallenge))

	//Replaced by enclave
	/*_, err := rand.Read(challenge)
	if err != nil {
		return nil, err
	}*/
	return challenge, nil
}

func (c Challenge) String() string {
	return base64.RawURLEncoding.EncodeToString(c)
}
