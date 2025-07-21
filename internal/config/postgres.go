package config

import "fmt"

type PostgresConfig struct {
	URL      string
	Host     string `yaml:"host" env:"POSTGRES_HOST"`
	Port     string `yaml:"port" env:"POSTGRES_PORT"`
	DB       string `yaml:"database" env:"POSTGRES_DATABASE"`
	User     string `yaml:"user" env:"POSTGRES_USER"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
	MaxConn  string `yaml:"max_conn" env:"POSTGRES_MAX_CONN"`
}

func (c *PostgresConfig) MakeURL() {
	c.URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_max_conns=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DB,
		c.MaxConn,
	)
}
