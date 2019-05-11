package main

import (
	"log"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/steam"
)

// TODO: callback url
func registerProviders() {
	i := 0
	if config.SteamKey != "" {
		i++
		steam := steam.New(config.SteamKey, "http://localhost:3000/steam/callback")
		goth.UseProviders(steam)
	}

	if config.DiscordKey != "" && config.DiscordSecret != "" {
		i++
		discord := discord.New(config.DiscordKey, config.DiscordSecret, "http://localhost:3000/discord/callback", discord.ScopeIdentify)
		goth.UseProviders(discord)
	}

	if i == 0 {
		log.Fatalf("no provivider key/secrets were providerd\n")
	}
}
