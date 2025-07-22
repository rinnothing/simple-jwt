package config

type AuthConfig struct {
	AccessKey      string `yaml:"access_key"`
	RefreshKey     string `yaml:"refresh_key" `
	RefreshHashKey string `yaml:"refresh_hash_key"`
}
