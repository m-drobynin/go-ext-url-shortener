package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
)

func GetAppConfigDefaults() AppConfig {
	localhost := "localhost"

	return AppConfig{
		Version: "0.0.0",
		NetAddress: NetAddress{
			Host: localhost,
			Port: 8080,
		},
		BaseURL: "http://" + localhost,
	}
}

func GetAppConfig() (*AppConfig, error) {
	config := GetAppConfigDefaults()

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Version %s; Usage of %s:\n", config.Version, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Var(&config.NetAddress, "a", "Net address host:port")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "base shortener url")
	flag.Parse()

	var envConfig EnvConfig
	err := env.Parse(&envConfig)

	if err != nil {
		return nil, err
	}

	// override with env values if exists
	err = config.SetFromEnv(envConfig)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
