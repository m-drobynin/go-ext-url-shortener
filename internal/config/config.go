package config

import (
	"flag"
	"fmt"
	"os"
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

func (address *NetAddress) String() string {
	return fmt.Sprintf("%s:%d", address.Host, address.Port)
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

func GetAppConfigDefaults() AppConfig {
	localhost := "http://localhost"

	return AppConfig{
		Version: "0.0.0",
		NetAddress: NetAddress{
			Host: localhost,
			Port: 8080,
		},
		BaseURL: localhost,
	}
}

func GetAppConfig() AppConfig {
	config := GetAppConfigDefaults()

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Version %s; Usage of %s:\n", config.Version, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Var(&config.NetAddress, "a", "Net address host:port")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "base shortener url")
	flag.Parse()

	return config
}
