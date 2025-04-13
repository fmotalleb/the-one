package notification

import "context"

type Notification struct {
	Ctx      context.Context
	Contacts []string
	Subject  string
	Message  string
}
