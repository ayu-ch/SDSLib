package types

import "time"

type AdminRequest string
type Status AdminRequest

const (
	NoRequest AdminRequest = "NoRequest"
	Pending   AdminRequest = "Pending"
	Accepted  AdminRequest = "Accepted"
	Denied    AdminRequest = "Denied"
)

type User struct {
	UserID       int
	Username     string
	Pass         string
	Role         string
	Created      time.Time
	AdminRequest AdminRequest
}

type BookRequests struct {
	RequestID   int
	UserID      int
	BookID      int
	RequestDate time.Time
	AcceptDate  time.Time
	Status      Status
}

type Books struct {
	BookID   int
	Title    string
	Author   string
	Genre    string
	Quantity int
}

type BooksWithStatus struct {
	Books
	Status Status
}

type BooksWithUser struct {
	Books
	BookRequests
	Username        string
	RequestDate     string
	RequestID       int
	UserID          int
	RequestedBookID int
	AvailableBookID int
}

type BookRequestsWithTitle struct {
	RequestID   int
	UserID      int
	BookID      int
	RequestDate time.Time
	Status      Status
	Title       string
}
