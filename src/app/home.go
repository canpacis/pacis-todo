package app

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
	. "github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/lucide"
	"github.com/canpacis/pacis/server"
	"github.com/canpacis/pacis/server/metadata"
	"github.com/canpacis/pacis/server/middleware"
	"github.com/canpacis/pacis/x"
)

type Home struct {
	Store *firestore.Client
}

func (*Home) Metadata() *metadata.Metadata {
	return &metadata.Metadata{
		Title: "Pacis Todo",
	}
}

func (h *Home) Page() Node {
	return Main(
		Class("bg-neutral-50 w-screen min-h-screen flex justify-center py-0 md:py-8"),

		Div(
			Class("md:max-w-md w-full h-fit p-6 flex flex-col gap-4"),

			Component(Profile),
			TodoForm(),
			server.Async(h.TodoList, P(Class("text-neutral-400 text-center my-8"), Text("Loading..."))),
		),
	)
}

func Profile(ctx context.Context) Node {
	user := middleware.GetUser[*auth.UserRecord](ctx)

	return Div(
		Attr("x-data", "auth"),
		Class("flex gap-2 items-center"),

		Span(
			Class("relative flex h-14 w-14 shrink-0 overflow-hidden rounded-full"),

			Img(Class("aspect-square h-full w-full"), Src(user.UserInfo.PhotoURL)),
		),
		Div(
			Class("flex flex-col gap-0"),

			P(Class("text-neutral-400 text-sm"), Text("Hello,")),
			P(Class("font-semibold"), Text(user.DisplayName)),
		),
		Button(
			Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 border border-input bg-white shadow-sm hover:bg-accent hover:text-accent-foreground h-9 px-3 border-neutral-300"),
			Class("ml-auto"),
			x.On("click", "logout()"),

			Text("Log out"),
		),
	)
}

func TodoForm() Node {
	return server.Form(
		"create",
		Class("flex flex-col gap-2"),

		Input(
			Required,
			Type("text"),
			Name("title"),
			Placeholder("Title"),
			Class("file:text-foreground border-neutral-300 placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground dark:bg-input/30 h-9 w-full min-w-0 rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-[color,box-shadow] outline-none file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm focus-visible:border-sky-400 focus-visible:ring-sky-400/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive"),
		),
		Textarea(
			Required,
			Name("description"),
			Placeholder("Description"),
			Class("border-neutral-300 placeholder:text-muted-foreground focus-visible:border-sky-400 focus-visible:ring-sky-400/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:bg-input/30 flex field-sizing-content min-h-16 w-full rounded-md border bg-transparent px-3 py-2 text-base shadow-xs transition-[color,box-shadow] outline-none focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50 md:text-sm"),
		),
		Div(
			Button(
				Type("submit"),
				Class("bg-sky-400 text-white shadow hover:bg-sky-400/90 inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-sky-400/50 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 h-9 px-3"),

				Text("Create"),
			),
		),
	)
}

type Todo struct {
	ID          string
	Title       string    `firestore:"title"`
	Description string    `firestore:"description"`
	Done        bool      `firestore:"done"`
	CreatedAt   time.Time `firestore:"created_at"`
}

func (t *Todo) Node() Node {
	return Div(
		Class("bg-white rounded-lg py-4 px-3 hover:shadow-lg shadow-neutral-100 flex gap-4 items-center"),

		Div(
			H1(Class("font-semibold text-lg"), Text(t.Title)),
			P(Class("text-neutral-500"), Text(t.Description)),
		),
		server.Form(
			"update",
			Class("ml-auto"),

			Input(Type("text"), Class("sr-only"), Name("id"), Value(t.ID)),
			Input(Type("text"), Class("sr-only"), Name("done"), If(t.Done, Value("false")), If(!t.Done, Value("true"))),
			Button(
				Type("submit"),
				Class("border border-neutral-300 rounded-full size-6 aspect-square shrink-0 flex items-center justify-center"),
				If(t.Done, Class("border-sky-400 bg-sky-400")),

				If(t.Done, lucide.Check(Class("size-4 text-white stroke-2"))),
				P(Class("sr-only"), Text("Complete")),
			),
		),
		server.Form(
			"delete",

			Input(Type("text"), Class("sr-only"), Name("id"), Value(t.ID)),
			Button(
				Type("submit"),
				Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 h-8 w-8 bg-red-500/10 text-red-500 hover:bg-red-500/20"),

				lucide.Trash2(),
			),
		),
	)
}

func GetTodos(ctx context.Context, store *firestore.Client) ([]Todo, error) {
	user := middleware.GetUser[*auth.UserRecord](ctx)
	collection := store.Collection("todos")
	query := collection.Where("owner", "==", user.UID).OrderBy("created_at", firestore.Desc).Documents(ctx)
	defer query.Stop()

	todos := []Todo{}

	for {
		doc, err := query.Next()
		if err != nil {
			break
		}
		var todo = Todo{
			ID: doc.Ref.ID,
		}
		if err := doc.DataTo(&todo); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (h *Home) TodoList(ctx context.Context) Node {
	todos, err := GetTodos(ctx, h.Store)
	if err != nil {
		return P(Textf("There was an error loading the todos: %s", err))
	}

	return Div(
		Class("flex flex-col gap-2"),

		Map(todos, func(todo Todo) Node {
			return todo.Node()
		}),
		If(len(todos) == 0, P(Class("text-neutral-400 text-center my-8"), Text("No Todos Yet"))),
	)
}

func (h *Home) Actions() map[string]server.ActionFunc {
	return map[string]server.ActionFunc{
		"create": h.Create,
		"update": h.Update,
		"delete": h.Delete,
	}
}

type CreateForm struct {
	Title       string `form:"title"`
	Description string `form:"description"`
}

func (h *Home) Create(w http.ResponseWriter, r *http.Request) error {
	form, err := server.FormData[CreateForm](r)
	if err != nil {
		return err
	}

	user := middleware.GetUser[*auth.UserRecord](r.Context())
	collection := h.Store.Collection("todos")
	_, _, err = collection.Add(r.Context(), map[string]any{
		"title":       form.Title,
		"description": form.Description,
		"done":        false,
		"owner":       user.UID,
		"created_at":  time.Now(),
	})
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

type UpdateForm struct {
	ID   string `form:"id"`
	Done bool   `form:"done"`
}

func (h *Home) Update(w http.ResponseWriter, r *http.Request) error {
	form, err := server.FormData[UpdateForm](r)
	if err != nil {
		return err
	}

	collection := h.Store.Collection("todos")
	_, err = collection.Doc(form.ID).Update(r.Context(), []firestore.Update{
		{Path: "done", Value: form.Done},
	})
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

type DeleteForm struct {
	ID string `form:"id"`
}

func (h *Home) Delete(w http.ResponseWriter, r *http.Request) error {
	form, err := server.FormData[DeleteForm](r)
	if err != nil {
		return err
	}

	collection := h.Store.Collection("todos")
	_, err = collection.Doc(form.ID).Delete(r.Context())
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}
