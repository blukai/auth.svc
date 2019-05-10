package main

import (
	"log"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/steam"
)

var providers = map[string]goth.Provider{}

func registerProviders() {
	// TODO: callback url

	if config.SteamKey != "" {
		steam := steam.New(config.SteamKey, "http://localhost:3000/steam/callback")
		goth.UseProviders(steam)
		providers[steam.Name()] = steam
	}

	if config.DiscordKey != "" && config.DiscordSecret != "" {
		discord := discord.New(config.DiscordKey, config.DiscordSecret, "http://localhost:3000/discord/callback", discord.ScopeIdentify)
		goth.UseProviders(discord)
		providers[discord.Name()] = discord
	}

	if len(providers) == 0 {
		log.Fatalf("no provivider key/secrets were providerd\n")
	}
}
