package app

import (
	"fmt"
	"mymail/app/pages"
	"mymail/app/partials"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/pocketbase/pocketbase/tools/hook"

	"github.com/alexedwards/scs/v2"
	"github.com/urfave/negroni"
)

var sessionManager *scs.SessionManager

const authCookieKey = "auth"

// About e.Auth:
// - func loadAuthToken() *hook.Handler[*core.RequestEvent]
// - https://github.com/pocketbase/pocketbase/blob/master/apis/middlewares.go#L181
func requireAuth(e *core.RequestEvent) (err error) {
	if e.Auth == nil {
		return e.Redirect(303, "/login")
	}
	return e.Next()
}

func htmxRedirect(e *core.RequestEvent, location string) {
	// redirection for users browsing the cliend
	e.Response.Header().Set("HX-Redirect", location)
}

// taking the pocketbase header-based implementation as reference
// - https://github.com/pocketbase/pocketbase/blob/e29655aba90817ed39d182a6b0f8056cdb15b069/apis/middlewares.go#L181
func loadAuthTokenFromCookie() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id:       "customLoadAuthToken",
		Priority: apis.DefaultLoadAuthTokenMiddlewarePriority - 1, // execute this auth middleware first
		Func: func(e *core.RequestEvent) error {
			// println("loadAuthTokenFromCookie")
			if e.Auth != nil {
				// already loaded by another middleware
				return e.Next()
			}

			token := sessionManager.GetString(e.Request.Context(), authCookieKey)

			if token == "" {
				return e.Next()
			}

			record, err := e.App.FindAuthRecordByToken(token, core.TokenTypeAuth)
			if err != nil {
				e.App.Logger().Debug("loadAuthToken failure", "error", err)
			} else if record != nil {
				e.Auth = record
			}

			return e.Next()
		},
	}
}

// does not to anything, yet, but useful for debugging
func rootMiddleware(e *core.RequestEvent) (err error) {

	lrw := negroni.NewResponseWriter(e.Response)
	e.Response = lrw
	// fmt.Println("rootMiddleware")
	err = e.Next()
	if err != nil && strings.Contains(fmt.Sprintf("%s", e.Request.URL), "auth-with-password") {
		// awkward code, but necessary to return Templ component inside
		// pocketbase middleware
		r := e.Request
		w := e.Response
		partials.MessageBox("Wrong credentials").Render(r.Context(), w)
		return e.HTML(http.StatusUnauthorized, "")
	}
	if err != nil {
		// return html error message here 'something went wrong, go back'
	}

	// [ debug ]
	// Useful to check the Response status after the route resolves:
	// fmt.Printf("<-- %d %s", lrw.Status(), http.StatusText(lrw.Status()))
	return err
}

func cookieMiddleware() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id:       "loadCookies",
		Priority: apis.DefaultLoadAuthTokenMiddlewarePriority - 2, // execute this auth middleware first
		Func:     apis.WrapStdMiddleware(sessionManager.LoadAndSave),
	}
}

func MakeRouter(app *pocketbase.PocketBase) {
	// cookies session manager
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	app.OnRecordAuthRequest().BindFunc(func(e *core.RecordAuthRequestEvent) error {
		// println("OnRecordAuthRequest")
		// to set cookie value
		sessionManager.Put(e.Request.Context(), authCookieKey, e.Token)

		// redirect successful logins
		e.Response.Header().Set("HX-Redirect", "/page/inbox")
		return e.NoContent(http.StatusOK)
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {

		se.Router.BindFunc(rootMiddleware)
		se.Router.Bind(cookieMiddleware())
		se.Router.Bind(loadAuthTokenFromCookie())

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

		se.Router.GET("/login",
			func(e *core.RequestEvent) error {
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if e.Auth != nil {
						// already logged in
						http.Redirect(w, r, "/page/inbox", http.StatusSeeOther)
					}
					pages.Login().Render(r.Context(), w)
				}).ServeHTTP(e.Response, e.Request)
				return nil
			})

		se.Router.POST("/logout",
			func(e *core.RequestEvent) error {
				sessionManager.Destroy(e.Request.Context())
				htmxRedirect(e, "/login")
				return e.NoContent(http.StatusOK)
			})

		pageGroup := se.Router.Group("/page")
		pageGroup.BindFunc(requireAuth) // require auth for /page* routes

		pageGroup.GET("/inbox", apis.WrapStdHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				pages.Inbox().Render(r.Context(), w)
			})))

		// CATCH-ALL
		se.Router.GET("/", func(e *core.RequestEvent) (err error) {
			return e.String(http.StatusNotFound, "This page does not exist, please go back... TODO: make friendly 404 page")
		})

		return se.Next()
	})
}
