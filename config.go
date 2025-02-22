package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ProxyHost   string
	ForwardHost string
}

var config *Config = nil

func GetConfig() *Config {
	if config == nil {
		loadConfig()
	}
	return config
}

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	config = &Config{
		ProxyHost:   getEnv("PROXY_HOST", "localhost:8888"),
		ForwardHost: getEnv("FORWARD_HOST", "localhost:3306"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
