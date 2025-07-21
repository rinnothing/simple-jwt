package config

type AuthConfig struct {
	AccessKey      string `yaml:"access_key" env:"ACCESS_KEY"`
	RefreshKey     string `yaml:"refresh_key" env:"REFRESH_KEY"`
	RefreshHashKey string `yaml:"refresh_hash_key" env:"REFRESH_HASH_KEY"`
}
