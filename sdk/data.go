package sdk

import "html/template"

// Event for module communication
type Event struct {
	// ID of the event
	ID string
	// Event name
	Event string
	// Data of the event
	Data interface{}
}

// MetaData contains the metadata of content
type MetaData struct {
	Title       string
	Description string
	Author      string
	Layouts     []string
	Scripts     []string
	Tags        []string
	Name        string
}

// SafeContent contains save HTML content
type SafeContent struct {
	MetaData
	Content template.HTML
}

// Content contains unsafe content
type Content struct {
	MetaData
	Content string
}

var (
	// ClientEvents queues events for the client
	ClientEvents chan Event
	// UpdateEvents queue for versioning
	UpdateEvents chan Event
)
