package router

import (
	"context"
	"mymail/app/pages"
	"mymail/app/shared"
	"net/http"

	"github.com/mrz1836/postmark"
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
					e.App.Logger().Error("record not found on delete:", "error", err)
					return
				}
				if record.GetString("to") != e.Auth.GetString("email") {
					http.Error(w, "", http.StatusBadRequest)
					e.App.Logger().Error("delete failed", "email", e.Auth.GetString("email"), "from", record.GetString("from"))
					return
				}
				record.Set("deleted", true)
				err = pb.Save(record)
				if err != nil {
					http.Error(w, "", http.StatusBadRequest)
					e.App.Logger().Error("could not delete:", "error", err)
					return
				}
			}).ServeHTTP(e.Response, e.Request)
			return nil
		})

	pageGroup.POST("/new/submit",
		func(e *core.RequestEvent) error {
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// println("/new/submit")
				collection, err := pb.FindCollectionByNameOrId("outbound")
				if err != nil {
					pages.Submitted("There was an error :(").Render(r.Context(), w)
					return
				}
				info, err := e.RequestInfo()
				if err != nil {
					e.App.Logger().Error("Failed RequestInfo")
					pages.Submitted("There was an error :(").Render(r.Context(), w)

				}
				to, ok := info.Body["to"].(string)
				subject, okk := info.Body["subject"].(string)
				body, okkk := info.Body["body"].(string)
				if !(ok && okk && okkk) {
					e.App.Logger().Error("Failed user input")
					pages.Submitted("There was an error :(").Render(r.Context(), w)
					return
				}

				record := core.NewRecord(collection)
				record.Set("to", to)
				record.Set("from", e.Auth.GetString("email"))
				record.Set("body", body)
				record.Set("subject", subject)

				error := pb.Validate(record)
				if error != nil {
					e.App.Logger().Error("Failed pbValidate")
					e.App.Logger().Error("pbValidate:", "error", error)
					pages.Submitted("There was an error :(").Render(r.Context(), w)

				}

				email := postmark.Email{
					From:     e.Auth.GetString("email"),
					To:       to,
					Subject:  subject,
					TextBody: body,
					Tag:      "outbound",
				}

				_, errApi := postmarkClient.SendEmail(context.Background(), email)
				if errApi != nil {
					e.App.Logger().Error("errApi:", "error", errApi)
					e.App.Logger().Error("errApi:", "error", errApi)
					e.App.Logger().Error("errApi:", "error", errApi)

					pages.Submitted("There was an error :(").Render(r.Context(), w)

					return
				}
				pb.Save(record)
				pages.Submitted("Your email has been sent!").Render(r.Context(), w)

			}).ServeHTTP(e.Response, e.Request)
			return nil
		})

	pageGroup.GET("/new",
		func(e *core.RequestEvent) error {
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				pages.New(e.Auth.GetString("email")).Render(r.Context(), w)
			}).ServeHTTP(e.Response, e.Request)
			return nil
		})

}
