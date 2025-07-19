package config

type AuthConfig struct {
	AccessKey      string `json:"access_key" env:"ACCESS_KEY"`
	RefreshKey     string `json:"refresh_key" env:"REFRESH_KEY"`
	RefreshHashKey string `json:"refresh_hash_key" env:"REFRESH_HASH_KEY"`
}
