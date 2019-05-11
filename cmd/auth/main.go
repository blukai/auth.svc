package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/blukai/auth.svc/pkg/handlers"
	"github.com/blukai/auth.svc/pkg/middleware"
	"github.com/blukai/auth.svc/pkg/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	version := flag.Bool("version", false, "show auth-svc version")
	flag.Parse()
	cmd := flag.Arg(0)

	if *version || cmd == "version" {
		fmt.Printf("version: x; repo: %s; commit: %s\n", Repo, Commit)
		return
	} else if cmd == "help" {
		// TODO: proper usage
		flag.Usage()
		return
	}

	loadConfig()
	registerProviders()

	store, err := redis.NewStore(10, "tcp", config.RedisAddress, "", []byte(config.SessionSecret))
	// TODO: proper session config, security
	if err != nil {
		log.Fatalf("could not create new redis store: %v\n", err)
	}

	router := gin.Default()
	router.Use(
		cors.New(cors.Config{
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return true
			},
		}),
		middleware.ProviderName(),
		sessions.Sessions("id", store),
	)
	router.GET("/:provider", handlers.Root)
	router.GET("/:provider/callback", handlers.ProviderCallback)
	router.POST("/logout", handlers.Logout)

	server := &http.Server{
		Handler: router,
		Addr:    config.ServiceHost + ":" + config.ServicePort,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("listen %v\n", err)
		}
	}()

	util.RegisterShutdown(5 * time.Second)
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("failed to shutdown server gracefully: %v\n", err)
	}
}
