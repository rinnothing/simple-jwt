package config

type PostgresConfig struct {
	KeysPassword string `json:"password" env:"KEYS_PASSWORD"`
}
