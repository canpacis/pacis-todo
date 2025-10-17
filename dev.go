//go:build dev

package main

import (
	"log/slog"
	"net/http"
	"net/url"

	"github.com/canpacis/pacis/server"
)

func init() {
	dev, _ := url.Parse("http://localhost:5173")
	options = &server.Options{
		Env:       server.Dev,
		Port:      ":8081",
		DevServer: dev,
		Logger:    slog.Default(),
		Mux:       http.NewServeMux(),
	}
}

func Ready(s *server.Server) {
	s.RegisterDevHandlers()
}
