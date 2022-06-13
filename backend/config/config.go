package config

type Config struct {
	Host   string `env:"SERVERHOST" validate:"required"`
	Port   uint16 `env:"SERVERPORT" validate:"required"`
	Infura InfuraConfig
}

type InfuraConfig struct {
	InfuraAddress string `env:"INFURA_ADDRESS" validate:"required"`
}
