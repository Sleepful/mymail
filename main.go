package main

import (
	// "bytes"
	"mymail/app"

	"fmt"
	"log"
	"os"
	"strings"

	// "github.com/a-h/templ"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
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
	pbInstance := pocketbase.NewWithConfig(pocketbase.Config{
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
	app.MakeRouter(pbInstance)

	// MIGRATIONS
	//
	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(pbInstance, pbInstance.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	// RUN
	//
	if err := pbInstance.Start(); err != nil {
		log.Fatal(err)
	}
}
