package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/DisposaBoy/JsonConfigReader"
)

type Config struct {
	Addr string
}

const (
	defaultAddr = ":8080"
)

func DefaultConfig() Config {
	config := Config{}
	config.Addr = ":8080"
	return config
}

func ReadConfig(filename string) (Config, error) {
	config := DefaultConfig()

	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	reader := JsonConfigReader.New(file)
	err = json.NewDecoder(reader).Decode(&config)
	if err != nil {
		return config, err
	}

	if config.Addr == "" {
		config.Addr = ":8080"
	}

	log.Printf("Config loaded from %q: %#v", filename, config)
	return config, nil
}
