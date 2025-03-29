package app

import (
	"mymail/app/pages"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/urfave/negroni"
)

// About e.Auth:
// - func loadAuthToken() *hook.Handler[*core.RequestEvent]
// - https://github.com/pocketbase/pocketbase/blob/master/apis/middlewares.go#L181
func requireAuth(e *core.RequestEvent) (err error) {
	if e.Auth == nil {
		return e.Redirect(301, "/login")
	}
	return e.Next()
}

// does not to anything, yet, but useful for debugging
func rootMiddleware(e *core.RequestEvent) (err error) {
	lrw := negroni.NewResponseWriter(e.Response)
	e.Response = lrw
	err = e.Next()
	// Useful to check the Response status after the route resolves:
	// fmt.Printf("<-- %d %s", lrw.Status(), http.StatusText(lrw.Status()))
	return err
}

func MakeRouter(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {

		pageGroup := se.Router.Group("/page")

		se.Router.BindFunc(rootMiddleware)

		pageGroup.BindFunc(requireAuth) // require auth for this route group

		se.Router.GET("/login", apis.WrapStdHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				pages.Layout(pages.Login()).Render(r.Context(), w)
			})))

		return se.Next()
	})
}
