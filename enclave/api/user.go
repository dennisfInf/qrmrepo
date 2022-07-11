package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

const StoragePath = "/data/user.json"

// User represents the user model
type User struct {
	id          uint64
	name        string
	displayName string
	credentials []webauthn.Credential
}

// NewUser creates and returns a new User
func NewUser(name, displayName string) User {
	return User{
		id:          randomUint64(),
		name:        name,
		displayName: displayName,
		credentials: make([]webauthn.Credential, 0),
	}
}

func randomUint64() uint64 {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return binary.LittleEndian.Uint64(buf)
}

// WebAuthnID returns the user's ID
func (u User) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, u.id)
	return buf
}

// WebAuthnName returns the user's username
func (u User) WebAuthnName() string {
	return u.name
}

// WebAuthnDisplayName returns the user's display name
func (u User) WebAuthnDisplayName() string {
	return u.displayName
}

// WebAuthnIcon is not (yet) implemented
func (u User) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user
func (u *User) AddCredential(cred webauthn.Credential) {
	u.credentials = append(u.credentials, cred)
}

// WebAuthnCredentials returns credentials owned by the user
func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

// CredentialExcludeList returns a CredentialDescriptor array filled
// with all the user's credentials
func (u User) CredentialExcludeList() []protocol.CredentialDescriptor {
	credentialExcludeList := make([]protocol.CredentialDescriptor, 0)

	for _, cred := range u.credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}

func (u User) startedRegistration() bool {
	return u.name != ""
}

func (u User) finishedRegistration() bool {
	return len(u.credentials) > 0
}

func SaveUser(u User) error {
	f, err := os.OpenFile(StoragePath, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(&u); err != nil {
		return fmt.Errorf("failed to write json to file: %w", err)
	}

	return nil
}

func LoadUser() (User, error) {
	f, err := os.Open(StoragePath)
	if os.IsNotExist(err) {
		// User does not exist yet
		return User{}, nil
	} else if err != nil {
		return User{}, err
	}
	defer f.Close()

	var user User
	if err := json.NewDecoder(f).Decode(&user); err != nil {
		return User{}, err
	}

	return user, nil
}
