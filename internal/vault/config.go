package vault

import "os"

type Config struct {
	Port    string
	Secret  string
	Storage string
}

func NewConfig() Config {
	return Config{
		Port:    os.Getenv("VAULT_PORT"),
		Secret:  os.Getenv("VAULT_SECRET"),
		Storage: os.Getenv("VAULT_STORAGE"),
	}
}
