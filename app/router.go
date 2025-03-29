package app

import (
	"fmt"
	"mymail/app/pages"
	"net/http"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	// "github.com/pocketbase/pocketbase/tools/hook"
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
	if err != nil && strings.Contains(fmt.Sprintf("%s", e.Request.URL), "auth-with-password") {
		fmt.Printf("<-- %s", err)
		// return html error message here 'wrong credentials'
	}
	if err != nil {
		// return html error message here 'something went wrong, go back'
	}

	// Useful to check the Response status after the route resolves:
	// fmt.Printf("<-- %d %s", lrw.Status(), http.StatusText(lrw.Status()))
	return err
}

func MakeRouter(app *pocketbase.PocketBase) {
	app.OnRecordAuthRequest().BindFunc(func(e *core.RecordAuthRequestEvent) error {
		// redirect successful logins
		return e.Redirect(301, "/site")
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {

		se.Router.BindFunc(rootMiddleware)

		// serve static assets
		se.Router.GET("/assets/{path...}", apis.Static(os.DirFS("assets"), false)).BindFunc(
			func(e *core.RequestEvent) error {
				if strings.Contains(fmt.Sprintf("%s", e.Request.URL), "index.js") {
					e.Response.Header().Set("Content-type", "text/javascript")
				}
				if strings.Contains(fmt.Sprintf("%s", e.Request.URL), "styles.css") {
					e.Response.Header().Set("Content-type", "text/css")
				}
				return e.Next()
			})

		pageGroup := se.Router.Group("/page")
		pageGroup.BindFunc(requireAuth) // require auth for /page* routes

		se.Router.GET("/login", apis.WrapStdHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				pages.Layout(pages.Login()).Render(r.Context(), w)
			})))

		/// ends OnServe
		//
		//
		return se.Next()
	})
}
