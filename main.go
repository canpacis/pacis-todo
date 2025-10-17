package main

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"github.com/canpacis/pacis-app/src/app"
	"github.com/canpacis/pacis/server"
	"github.com/canpacis/pacis/server/middleware"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var options *server.Options

func main() {
	godotenv.Load()

	fbapp, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(credentials))
	if err != nil {
		log.Fatal(err)
	}
	firestore, err := fbapp.Firestore(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	server := server.New(options)

	Ready(server)

	auth := middleware.NewAuthentication(&app.Authenticator{App: fbapp})
	home := &app.Home{Store: firestore}
	server.HandlePage("/", home, app.RootLayout, middleware.DefaultGzip, auth)
	server.HandlePage("GET /login", &app.Login{}, app.RootLayout, middleware.DefaultGzip)
	server.HandleFunc("GET /login/done", app.LoginDone)
	server.HandleFunc("GET /logout", app.Logout)

	server.Serve()
}
