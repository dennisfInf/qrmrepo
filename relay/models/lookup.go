package models

type Lookup struct {
	ID             uint64 `db:"id"`
	Username       string `db:"username"`
	EnclaveAddress string `db:"enclave_address"`
}
