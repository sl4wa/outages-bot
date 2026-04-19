package notifier

import "time"

// Content carries the structured data needed to render an outage notification.
type Content struct {
	City       string
	StreetName string
	Buildings  []string
	Start      time.Time
	End        time.Time
	Comment    string
}

// Sender sends notifications to users.
type Sender interface {
	Send(userID int64, content Content) error
}
