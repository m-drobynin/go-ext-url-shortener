package config

import (
	"fmt"
	"strconv"
	"strings"
)

type AppConfig struct {
	Version    string
	NetAddress NetAddress
	BaseURL    string
}

type NetAddress struct {
	Host string
	Port int
}

type EnvConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func (address *NetAddress) String() string {
	return fmt.Sprintf("http://%s:%d", address.Host, address.Port)
}

func (address *NetAddress) Set(flagValue string) error {
	parts := strings.Split(flagValue, ":")
	address.Host = parts[0]

	if len(parts) == 2 {
		port, err := strconv.Atoi(parts[1])

		if err != nil {
			return err
		}

		address.Port = port
	}

	return nil
}

func (config *AppConfig) String() string {
	return fmt.Sprintf("Version: %s, address: %s:%d", config.Version, config.NetAddress.Host, config.NetAddress.Port)
}

func (config *AppConfig) GetBaseURL() string {
	return config.NetAddress.String()
}

func (config *AppConfig) SetFromEnv(envConfig EnvConfig) error {
	if envConfig.BaseURL != "" {
		config.BaseURL = envConfig.BaseURL
	}

	if envConfig.ServerAddress != "" {
		err := config.NetAddress.Set(envConfig.ServerAddress)

		if err != nil {
			return err
		}
	}

	return nil
}
