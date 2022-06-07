package config

type ServerConfig struct {
	Host           string `env:"SERVERHOST" validate:"required"`
	Port           uint16 `env:"SERVERPORT" validate:"required"`
	KubeConfigPath string `env:"KUBECONFIGPATH" validate:"required"`
}

type PostgresConfig struct {
	Host     string `env:"DBHOST" validate:"required"`
	Port     uint16 `env:"DBPORT" validate:"required"`
	User     string `env:"DBUSER" validate:"required"`
	Password string `env:"DBPASSWORD" validate:"required"`
	DBName   string `env:"DBNAME" validate:"required"`
}

type GlobalConfig struct {
	Server   ServerConfig
	Postgres PostgresConfig
}
