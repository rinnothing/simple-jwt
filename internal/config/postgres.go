package config

import "fmt"

type PostgresConfig struct {
	URL      string
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DB       string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	MaxConn  string `yaml:"max_conn"`
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
