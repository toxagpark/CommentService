package entity

import "time"

type Comment struct {
	ID        int
	User      User
	Text      string
	CreatedAt time.Time
	IsEdited  bool
}
