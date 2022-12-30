package s

import (
	"html/template"
	"time"
)

var Temp *template.Template

type StatusData struct {
	StatusCode int
	StatusMsg  string
}

// Create a struct that models the structure of a user in the request body
type Credentials struct {
	Password string
	Username string
}

var sessions = map[string]session{}

// each session contains the username of the user and the time at which it expires
type session struct {
	username string
	expiry   time.Time
}
type User struct {
	// not sure about userId!
	UserId   int
	Username string
	Pass     string
	Email    string
	Time     string
}
type Post struct {
	PostId       int
	UserId       int
	Title        string
	Content      string
	CreationTime string
}

// for testing
type PostTopic struct {
	PostId  int
	TopicId int
}
type Comment struct {
	UserId  int
	PostId  int
	Content string
	Time    string
}
type Reaction struct {
	ReactionId int
	UserId     int
	PostId     int
	CommentId  int
	Reaction   string
}
type Colours struct {
	Reset      string // Resets terminal colour to default after 'text colouring'
	Red        string
	LightRed   string
	Green      string
	LightGreen string
	Blue       string
	LightBlue  string
	Orange     string
	Yellow     string
}
