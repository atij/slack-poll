package app

import (
	"github.com/atij/slack-poll/model"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
)

type payload struct {
	Data string `form:"payload"`
}

// Help response structure
func getHelpReponse(c *gin.Context) *slack.Message {
	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "*Adore Me poll for slack*", false, false),
			nil,
			nil,
		),
		slack.NewDividerBlock(),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "*Simple poll*", false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "\n/am-poll \"What's your favourite color ?\" \"Red\" \"Green\" \"Blue\" \"Yellow\"\n", false, false),
			nil,
			nil,
		),
	}

	msg := slack.NewBlockMessage(blocks...)

	return &msg
}

// Poll response ...
func getPollResponse(p *model.Poll) *slack.Message {

	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", p.Title, false, false),
			nil,
			nil,
		),
		slack.NewDividerBlock(),
	}

	for _, item := range p.Options {
		be := slack.NewButtonBlockElement(p.ID + "::" + item.Title, item.Title, slack.NewTextBlockObject("plain_text", item.Title, false, false))
		bk := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", item.Title, false, false), nil, slack.NewAccessory(be))
		
		blocks = append(blocks, bk)

		var voters []slack.MixedElement
		for _, v := range item.Votes {
			voters = append(voters, *slack.NewImageBlockElement(v.Avatar, v.UserName))
		}
		if len(voters) != 0 {
			ct := slack.NewContextBlock(item.Title + " voters", voters...)
			blocks = append(blocks, ct)
		}

	}

	blocks = append(blocks, slack.NewDividerBlock())
	ab := slack.NewActionBlock(p.ID + "::actions", slack.NewButtonBlockElement(p.ID+ "::addOption", "Add Option", slack.NewTextBlockObject("plain_text", "Add Option", true, false)))
	blocks = append(blocks, ab)

	msg := slack.NewBlockMessage(blocks...)
	msg.Msg.ReplaceOriginal = true
	msg.Msg.ResponseType = slack.ResponseTypeInChannel

	return &msg
}

