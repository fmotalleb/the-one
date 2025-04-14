package notification

import "context"

type Notification struct {
	Ctx           context.Context
	ContactPoints []string
	Subject       string
	Message       string
}
