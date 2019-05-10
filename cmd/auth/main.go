package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/blukai/auth.svc/pkg/handlers"
	"github.com/blukai/auth.svc/pkg/middleware"
	"github.com/blukai/auth.svc/pkg/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/steam"
)

func main() {
	goth.UseProviders(
		steam.New(util.GetEnvOrDie("STEAM_KEY"), "http://localhost:3000/steam/callback"),
		discord.New(os.Getenv("DISCORD_KEY"), os.Getenv("DISCORD_SECRET"), "http://localhost:3000/discord/callback", discord.ScopeIdentify),
	)

	store, err := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		log.Fatalf("could not create new redis store: %v", err)
	}

	router := gin.Default()

	router.Use(
		sessions.Sessions("id", store),
		middleware.ProviderName(),
	)

	router.GET(":provider", handlers.Provider())
	router.GET(":provider/callback", handlers.ProviderCallback())
	router.GET(":provider/logout", handlers.ProviderLogout())
	router.GET(":provider/user", handlers.ProviderUser())

	server := &http.Server{Handler: router, Addr: ":3000"}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("listen %v\n", err)
		}
	}()

	util.RegisterShutdown(5 * time.Second)
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("failed to shutdown server gracefully: %v", err)
	}
}
