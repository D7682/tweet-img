package config

type Config struct {
	ApiKey       string `yaml:"api_key"`
	ApiSecret    string `yaml:"api_secret"`
	BearerToken  string `yaml:"bearer_token"`
	AccessToken  string `yaml:"access_token"`
	AccessSecret string `yaml:"access_secret"`
}
