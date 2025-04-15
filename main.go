package main

import (
	"log/slog"
	"net/http"
	"os"

	_ "net/http/pprof"

	_ "github.com/joho/godotenv/autoload"

	"github.com/adriancmiranda-wm/cli/cmd"
	"github.com/adriancmiranda-wm/cli/internal/log"
)

// Variáveis preenchidas pelo ldflags do GoReleaser
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	defer log.RecoverPanic("main", func(err error) {
		slog.Error("Application terminated due to unhandled panic", "error", err)
	})

	if os.Getenv("WM_PROFILE") != "" {
		go func() {
			slog.Info("Serving pprof at localhost:6060")
			if httpErr := http.ListenAndServe("localhost:6060", nil); httpErr != nil {
				slog.Error("Failed to serve pprof", "error", httpErr)
			}
		}()
	}

	slog.Info("WM CLI starting",
		"version", Version,
		"commit", Commit,
		"date", Date,
	)

	cmd.Execute()
}
