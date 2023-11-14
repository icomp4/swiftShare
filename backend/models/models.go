package models

import "time"

type User struct {
	ID        uint
	Username  string
	Email     string
	Password  string
	PfpUrl    string
	Files     []File
	CreatedAt time.Time
	DeletedAt time.Time
}
type File struct {
	ID        uint
	UserID    uint
	Name      string
	Path      string
	Link      string
	CreatedAt time.Time
	DeletedAt time.Time
}