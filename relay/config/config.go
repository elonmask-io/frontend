package config

type ServerConfig struct {
	Port uint16 `env:"SERVERPORT" validate:"required"`
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
