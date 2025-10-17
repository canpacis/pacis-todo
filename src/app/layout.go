package app

import (
	. "github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/server"
	"github.com/canpacis/pacis/server/font"
)

func RootLayout(server *server.Server, head Node, children Node) Node {
	return Fragment(
		Doctype,
		Head(
			Meta(Charset("UTF-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1.0")),

			head,
			Link(Rel("icon"), Type("image/x-icon"), Href("/favicon.ico")),
			Link(Rel("stylesheet"), Href(server.Asset("/src/web/style.css"))),
			Script(Defer, Type("module"), Src(server.Asset("/src/web/main.ts"))),
			font.Head(
				font.New("Inter", font.WeightList{font.W100, font.W900}, font.Auto, font.Latin, font.LatinExt),
			),
		),
		Body(
			children,
			Script(Src(server.Asset("/src/web/stream.ts"))),
		),
	)
}
