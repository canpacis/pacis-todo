//go:build dev

package main

import (
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/canpacis/pacis/server"
)

var credentials []byte

func init() {
	dev, _ := url.Parse("http://localhost:5173")
	options = &server.Options{
		Env:       server.Dev,
		Port:      ":8081",
		DevServer: dev,
		Logger:    slog.Default(),
		Mux:       http.NewServeMux(),
	}
	var err error
	credentials, err = os.ReadFile("account-key.json")
	if err != nil {
		log.Fatal(err)
	}
}

func Ready(s *server.Server) {
	s.RegisterDevHandlers()
}
