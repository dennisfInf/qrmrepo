package persistence

var schemas = [...]string{
	`CREATE TABLE IF NOT EXISTS Lookup (
		id BIGSERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		enclave_address TEXT UNIQUE NOT NULL
	);`,
}
