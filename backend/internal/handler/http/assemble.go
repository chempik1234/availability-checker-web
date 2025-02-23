package http

import (
	"github.com/chempik1234/availability-checker-web/internal/ports/logs/logsadapters"
	"github.com/chempik1234/availability-checker-web/internal/ports/tokens/tokensadapters"
	"net/http"
)

func Assemble(
	logsRepo *logsadapters.LogRecordRepositoryDB,
	tokensRepo *tokensadapters.TokensRepositoryRedis,
) http.Handler {

	logsHandler := NewLogsHttpHandler(logsRepo)
	tokensHandler := NewTokensHttpHandler(tokensRepo)

	protectedWithEternalTokenMux := http.NewServeMux()
	protectedWithEternalTokenMux.HandleFunc("/", logsHandler.NewReceiveLogsHandler())
	protectedWithEternalTokenHandler := tokensHandler.CheckTokenMiddleware(protectedWithEternalTokenMux)

	siteHandler := http.NewServeMux()
	siteHandler.HandleFunc("/logs", logsHandler.NewLogsHandler())
	siteHandler.HandleFunc("/tokens", tokensHandler.NewTokensHandler())
	siteHandler.Handle("/", protectedWithEternalTokenHandler)

	return siteHandler
}
