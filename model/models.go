package model

import "errors"

// Poll struct  ...
type Poll struct {
	ID      string
	Text    string
	Owner   Owner
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
	Avatar   string
}

// Owner ...
type Owner struct {
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

// HasVote ...
func (p *Poll) HasVote(pollOption string, v Vote) bool {
	for _, item := range p.Options {
		if item.Title == pollOption {
			for _, vote := range item.Votes {
				if vote.UserID == v.UserID {
					return true
				}
			}
		}
	}
	return false
}

// RemoveVote ...
func (p *Poll) RemoveVote(pollOption string, v Vote) {
	for i, item := range p.Options {
		if item.Title == pollOption {
			for j, vote := range item.Votes {
				if vote.UserID == v.UserID {
					p.Options[i].Votes = append(p.Options[i].Votes[:j], p.Options[i].Votes[j+1:]...)
				}
			}
		}
	}
}