//go:build prod

package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"

	"github.com/canpacis/pacis/server"
)

func init() {
	options = &server.Options{
		Env:    server.Prod,
		Port:   ":8080",
		Logger: slog.Default(),
		Mux:    http.NewServeMux(),
	}
}

var (
	//go:embed build/*
	build embed.FS
	//go:embed build/static/.vite
	vite embed.FS
	//go:embed account-key.json
	credentials []byte
)

func Ready(s *server.Server) {
	if err := s.SetBuildDir("build/static", build, vite); err != nil {
		log.Fatal(err)
	}
}
