package main

import (
	"context"
	"fmt"
	"github.com/chempik1234/availability-checker-web/config"
	httpHandlers "github.com/chempik1234/availability-checker-web/internal/handler/http"
	"github.com/chempik1234/availability-checker-web/internal/ports/logs/logsadapters"
	"github.com/chempik1234/availability-checker-web/internal/ports/tokens/tokensadapters"
	"github.com/chempik1234/availability-checker-web/pkg/storage/postgres"
	"github.com/chempik1234/availability-checker-web/pkg/storage/redis"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	mainConfig, err := config.FromEnv()
	if err != nil {
		log.Fatal(fmt.Errorf("error loading config: %w", err))
	}

	DB, err := postgres.NewDBInstance(ctx, mainConfig.DB)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to relational database: %w", err))
	}
	RedisDB, err := redis.NewRedisClient(ctx, mainConfig.Redis.URL)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to Redis database: %w", err))
	}

	logsRepo := logsadapters.NewLogRecordRepositoryDB(DB)
	tokensRepo := tokensadapters.NewTokensRepositoryRedis(RedisDB)

	siteHandler := httpHandlers.Assemble(logsRepo, tokensRepo)

	log.Printf("serving at: %s\n", mainConfig.HTTP.Port)
	err = http.ListenAndServe(":"+mainConfig.HTTP.Port, siteHandler)
	if err != nil {
		log.Fatalf("error serving http server: %v", err)
	}
}
