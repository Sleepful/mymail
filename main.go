package main

import (
	"fmt"
	"log"
	"mymail/app/pages"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

type GlobalState struct {
	Count int
}

var global GlobalState
var sessionManager *scs.SessionManager

func getHandler(w http.ResponseWriter, r *http.Request) {
	userCount := sessionManager.GetInt(r.Context(), "count")
	component := pages.Page(global.Count, userCount)
	component.Render(r.Context(), w)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// Update state.
	r.ParseForm()

	// Check to see if the global button was pressed.
	if r.Form.Has("global") {
		global.Count++
	}
	if r.Form.Has("user") {
		currentCount := sessionManager.GetInt(r.Context(), "count")
		sessionManager.Put(r.Context(), "count", currentCount+1)
	}

	// Display the form.
	getHandler(w, r)
}

var dev = true

func disableCacheInDevMode(next http.Handler) http.Handler {
	if !dev {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize the session.
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	mux := http.NewServeMux()

	// Handle POST and GET requests.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			postHandler(w, r)
			return
		}
		getHandler(w, r)
	})

	mux.Handle("/assets/",
		disableCacheInDevMode(
			http.StripPrefix("/assets",
				http.FileServer(http.Dir("assets")))))

	// Add the middleware.
	muxWithSessionMiddleware := sessionManager.LoadAndSave(mux)

	// Start the server.
	fmt.Println("listening on http://localhost:8090")
	// use 0.0.0.0 instead of localhost for docker to resolve properly
	// https://stackoverflow.com/questions/64666842/docker-run-connection-was-reset-while-the-page-was-loading
	if err := http.ListenAndServe("0.0.0.0:8090", muxWithSessionMiddleware); err != nil {
		log.Printf("error listening: %v", err)
	}
}
