package shared

import (
	"context"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type contextKey string

var EmailListKey contextKey = "emailListKey"

type EmailData struct {
	Id      string
	From    string
	To      string
	Subject string
	Body    string
	Date    types.DateTime
}

func CreateEmailData(r *core.Record) EmailData {
	return EmailData{
		Id:      r.GetString("id"),
		From:    r.GetString("from"),
		To:      r.GetString("to"),
		Subject: r.GetString("subject"),
		Body:    r.GetString("body"),
		Date:    r.GetDateTime("created"),
	}

}

func GetEmailList(ctx context.Context) []EmailData {
	if list, ok := ctx.Value(EmailListKey).([]EmailData); ok {
		// fmt.Printf("list: %v\n", list)
		return list
	}
	return nil
}
