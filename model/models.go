package model

import "errors"

// Poll struct  ...
type Poll struct {
	ID      string
	Text    string
	Owner   string
	Channel string
	Title   string
	Options []PollOption
}

// PollOption ...
type PollOption struct {
	Title string
	Votes []Vote
}

// Vote ...
type Vote struct {
	UserID   string
	UserName string
}

// AddVote ...
func (p *Poll) AddVote(pollOption string, v Vote) error {
	for i, item := range p.Options {
		if item.Title == pollOption {
			p.Options[i].Votes = append(p.Options[i].Votes, v)
			return nil
		}
	}
	return errors.New("no poll option found")
}