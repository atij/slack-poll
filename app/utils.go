package app

import (
	"os"
	"strings"

	"github.com/atij/slack-poll/model"
	"github.com/slack-go/slack"
)

func cleanDoubleQuotes(i string) string {
	return strings.ReplaceAll(strings.ReplaceAll(i, "\u201D", "\""), "\u201C", "\"")
}

func splitOptions(s string) []string {
	var res []string
	// split options, remove quotes
	options := strings.Split(s, "\" " )
	for _, item := range options {
		res = append(res, strings.ReplaceAll(item, "\"", ""))
	}
	return res
}

func createPoll(c *slack.SlashCommand, options []string) (*model.Poll, error) {
	
	var opts []model.PollOption
	for _, item := range options[1:] {
		opts = append(opts, model.PollOption{
			Title: item,
		})
	}

	return &model.Poll{
		Text: options[0],
		Channel: c.ChannelName,
		Owner: c.UserName,
		Title: options[0],
		Options: opts,
	}, nil

}

func getUser(id string) (*slack.UserProfile, error) {
	client := slack.New(os.Getenv("SLACK_TOKEN"))
	profile, err := client.GetUserProfile(id, false)
	if err != nil {
		return nil, err
	}

	return profile, nil
}