package main

import (
	"embed"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"
)

// Die Oberfläche wird beim Kompilieren in das Programm aufgenommen.
// Andere Projektdateien wie Konfiguration oder Go-Code werden nicht ausgeliefert.
//
//go:embed index.html styles.css dist/frontend/main.js api/demo-prices.json diagrams/*.svg diagrams/*.d2
var assets embed.FS

func main() {
	if len(os.Args) > 2 || (len(os.Args) == 2 && os.Args[1] != "serve") {
		fmt.Fprintln(os.Stderr, "Aufruf: strom-dashboard [serve]")
		os.Exit(2)
	}

	addr := os.Getenv("ADDR")
	if addr == "" {
		port := os.Getenv("PORT")
		if port == "" {
			port = "4176"
		}
		addr = net.JoinHostPort("127.0.0.1", port)
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           newHandler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("Webserver kann nicht starten", "error", err)
		os.Exit(1)
	}
	fmt.Printf("Listening http://%s (Go)\n", listener.Addr())
	if err := server.Serve(listener); err != nil {
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
	mux.Handle("GET /assets/main.js", embeddedFile("dist/frontend/main.js", "text/javascript; charset=utf-8"))
	mux.Handle("GET /api/prices", embeddedFile("api/demo-prices.json", "application/json"))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok"}`)
	})

	return mux
}

func embeddedFile(path, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := assets.ReadFile(path)
		if err != nil {
			http.Error(w, "Datei nicht verfügbar", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", fmt.Sprint(len(data)))
		if r.Method != http.MethodHead {
			w.Write(data)
		}
	}
}
