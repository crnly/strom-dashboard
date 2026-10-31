package main

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Die Oberfläche wird beim Kompilieren in das Programm aufgenommen.
// Andere Projektdateien wie Konfiguration oder Go-Code werden nicht ausgeliefert.
//
//go:embed index.html styles.css diagrams/*.svg diagrams/*.d2
var assets embed.FS

func main() {
	if len(os.Args) > 2 || (len(os.Args) == 2 && os.Args[1] != "serve") {
		fmt.Fprintln(os.Stderr, "Aufruf: strom-dashboard [serve]")
		os.Exit(2)
	}

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:4176"
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           newHandler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("Webserver startet", "address", addr)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Webserver beendet", "error", err)
		os.Exit(1)
	}
}

func newHandler() http.Handler {
	mux := http.NewServeMux()
	files := http.FileServerFS(assets)

	// Nur diese öffentlichen Pfade erreichen den eingebetteten Dateiserver.
	mux.Handle("GET /{$}", files)
	mux.Handle("GET /styles.css", files)
	mux.Handle("GET /diagrams/{filename}", files)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok"}`)
	})

	return mux
}
