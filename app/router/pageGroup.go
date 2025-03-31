package router

import (
	"context"
	"mymail/app/pages"
	"mymail/app/shared"
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func getEmails(pb *pocketbase.PocketBase, e *core.RequestEvent) (context.Context, error) {
	records, err := pb.FindRecordsByFilter("inbound",
		"deleted = false && to = {:email}",
		"-created", 0, 0,
		dbx.Params{"email": e.Auth.GetString("email")})
	if err != nil {
		return nil, err
	}
	all := make([]shared.EmailData, len(records))
	for i, v := range records {
		data := shared.CreateEmailData(v)
		all[i] = data
		// fmt.Printf("%#v\n", data)
	}
	return context.WithValue(context.Background(), shared.EmailListKey, all), nil
}

func MakePageGroup(pb *pocketbase.PocketBase, se *core.ServeEvent) {

	pageGroup := se.Router.Group("/page")
	pageGroup.BindFunc(requireAuth) // require auth for /page* routes

	pageGroup.GET("/inbox",
		func(e *core.RequestEvent) error {
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Create a context variable that inherits from a parent, and sets the value "test".
				// Create a context key for the theme.
				ctx, err := getEmails(pb, e)
				if err != nil {
					http.Error(w, "", http.StatusBadRequest)
					return
				}
				// identify if it is an htmx request or not, in order to determine
				// if the layout ought to be included in the response
				hx := e.Request.Header.Get("HX-Request")
				showLayout := hx == ""
				pages.Inbox(e.Auth.GetString("email"), showLayout).Render(ctx, w)
			}).ServeHTTP(e.Response, e.Request)
			return nil
		})

	pageGroup.DELETE("/inbox/delete/{id}",
		func(e *core.RequestEvent) error {
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				record, err := pb.FindRecordById("inbound", id)
				if err != nil {
					http.Error(w, "", http.StatusBadRequest)
					return
				}
				if record.GetString("from") != e.Auth.GetString("email") {
					http.Error(w, "", http.StatusBadRequest)
					return
				}
				record.Set("deleted", true)
				err = pb.Save(record)
			}).ServeHTTP(e.Response, e.Request)
			return nil
		})

	pageGroup.GET("/new",
		func(e *core.RequestEvent) error {
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Create a context variable that inherits from a parent, and sets the value "test".
				// Create a context key for the theme.
				pages.New(e.Auth.GetString("email")).Render(r.Context(), w)
			}).ServeHTTP(e.Response, e.Request)
			return nil
		})

}

// func(e *core.RequestEvent) error {
// 	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// if record.GetString("from") != e.Auth.GetString("email") {
// 		// 	http.Error(w, "", http.StatusBadRequest)
// 		// 	return
// 		// }
// 	}).ServeHTTP(e.Response, e.Request)
// 	return nil
// })
