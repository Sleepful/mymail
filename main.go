package main

import (
	// "bytes"
	"mymail/app/pages"

	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	// "github.com/a-h/templ"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// register the libsql driver to use the same query builder
// implementation as the already existing sqlite3 builder
func init() {
	dbx.BuilderFuncMap["libsql"] = dbx.BuilderFuncMap["sqlite3"]
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// DATABASE
	//
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DBConnect: func(dbPath string) (*dbx.DB, error) {
			// turn off for first-time db creation to avoid turso timeout:
			// if false {
			if strings.Contains(dbPath, "data.db") {
				tursoToken := os.Getenv("TURSO_TOKEN")
				tursoUrl := os.Getenv("TURSO_URL")
				urlWithToken := fmt.Sprintf("%s?authToken=%s", tursoUrl, tursoToken)
				// fmt.Println(urlWithToken)
				return dbx.Open("libsql", urlWithToken)
			}

			// optionally for the logs (aka. pb_data/auxiliary.db) use the default local filesystem driver
			return core.DefaultDBConnect(dbPath)
		},
	})

	// ROUTES
	//
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {

		se.Router.BindFunc(apis.WrapStdMiddleware(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h.ServeHTTP(w, r)
			})
		}))
		// g := se.Router.Group("/page")
		// attach group middleware
		se.Router.BindFunc(func(e *core.RequestEvent) error {
			e.Next()
			return e.Next()
		})

		se.Router.GET("/page/{name}", apis.WrapStdHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// wraps component in layout
				pages.Layout(pages.Login()).Render(r.Context(), w)
			})))

		// register "GET /hello/{name}" route (allowed for everyone)
		se.Router.GET("/hello/{name}", func(e *core.RequestEvent) error {
			name := e.Request.PathValue("name")

			return e.String(http.StatusOK, "Hello "+name)
		})

		se.Router.GET("/page/login", func(e *core.RequestEvent) error {
			r := e.Request
			w := e.Response
			pages.Login().Render(r.Context(), w)
			// w := new(bytes.Buffer)
			// pages.Login().Render(r.Context(), w)
			// s := w.String()

			println(2)
			return e.HTML(http.StatusOK, "")
		})

		// register "POST /api/myapp/settings" route (allowed only for authenticated users)
		se.Router.POST("/api/myapp/settings", func(e *core.RequestEvent) error {
			// do something ...
			return e.JSON(http.StatusOK, map[string]bool{"success": true})
		}).Bind(apis.RequireAuth())

		return se.Next()
	})

	// MIGRATIONS
	//
	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	// RUN
	//
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
