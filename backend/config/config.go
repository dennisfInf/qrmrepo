package config

type Config struct {
	Port   uint16 `env:"SERVERPORT" validate:"required"`
	Infura InfuraConfig
}

type InfuraConfig struct {
	InfuraAddress string `env:"INFURAADDRESS" validate:"required"`
}
