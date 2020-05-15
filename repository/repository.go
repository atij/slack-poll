package repository

import "github.com/atij/slack-poll/model"

// Poll ...
type Poll interface {
	Create(p model.Poll) bool
	Find(string) (*model.Poll, error)
	Update(string, *model.Poll) error
}
