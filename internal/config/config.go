package config

import (
	"fmt"
	"os"
)

type Config struct {
	DSN    string
	UseK8s bool
	Debug  bool
}

func LoadConfigFromEnv() (Config, error) {
	_, ok := os.LookupEnv("DEBUG")
	conf := Config{
		DSN:    os.Getenv("DSN"),
		UseK8s: os.Getenv("USE_K8S") == "true",
		Debug:  ok,
	}
	if conf.DSN == "" {
		return Config{}, fmt.Errorf("missing env: %+v", conf)
	}
	return conf, nil
}
