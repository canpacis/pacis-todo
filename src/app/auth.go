package app

import (
	"fmt"
	"net/http"
	"time"

	firebase "firebase.google.com/go/v4"
	. "github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/lucide"
	"github.com/canpacis/pacis/server/metadata"
	"github.com/canpacis/pacis/x"
)

type Login struct{}

func (*Login) Metadata() *metadata.Metadata {
	return &metadata.Metadata{
		Title: "Login",
	}
}

func (*Login) Page() Node {
	return Div(
		Class("w-screen h-screen flex flex-col gap-2 items-center justify-center"),
		Attr("x-data", "auth"),

		lucide.GalleryVerticalEnd(Class("size-6")),
		H1(
			Class("font-semibold text-xl mb-4"),

			Text("Welcome to Pacis Todos"),
		),
		Button(
			Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 border border-input bg-white shadow-sm hover:bg-accent hover:text-accent-foreground h-9 px-3 border-neutral-300"),
			x.On("click", "login()"),

			Text("Login with Google"),
		),
	)
}

func LoginDone(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if len(token) == 0 {
		w.Write([]byte("Invalid token"))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Path:     "/",
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Hour),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Path:  "/",
		Name:  "token",
		Value: "",
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}

type Authenticator struct {
	App *firebase.App
}

func (a *Authenticator) Authenticate(r *http.Request) (any, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, fmt.Errorf("no cookie: %w", err)
	}

	auth, err := a.App.Auth(r.Context())
	if err != nil {
		return nil, fmt.Errorf("failed to get auth client: %w", err)
	}
	token, err := auth.VerifyIDToken(r.Context(), cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to verify id token: %w", err)
	}
	user, err := auth.GetUser(r.Context(), token.UID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (*Authenticator) OnError(w http.ResponseWriter, r *http.Request, err error) {
	http.Redirect(w, r, "/login", http.StatusFound)
}
