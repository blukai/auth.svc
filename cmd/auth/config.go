package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

var config struct {
	// service
	ServiceHost   string `default:"0.0.0.0" envconfig:"SERVICE_HOST"`
	ServicePort   string `default:"3000" envconfig:"SERVICE_PORT"`
	SessionSecret string `required:"true" envconfig:"SESSION_SECRET"`
	RedisAddress  string `default:"localhost:6379" envconfig:"REDIS_ADDRESS"`

	// providers
	SteamKey      string `envconfig:"STEAM_KEY"`
	DiscordKey    string `envconfig:"DISCORD_KEY"`
	DiscordSecret string `envconfig:"DISCORD_SECRET"`
}

func loadConfig() {
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("failed to load env config: %v\n", err)
	}
}
