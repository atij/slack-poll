package app

import (
	"fmt"

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
		slack.NewDividerBlock(),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "*Anonymous poll*", false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "\n/am-poll \"What's your favourite color ?\" \"Red\" \"Green\" \"Blue\" \"Yellow\" anonymous\n", false, false),
			nil,
			nil,
		),
	}

	msg := slack.NewBlockMessage(blocks...)

	return &msg
}

// Poll response ...
func getPollResponse(p *model.Poll) *slack.Message {

	var blocks []slack.Block

	var txt string
	if p.Mode.Anonymous {
		txt = "*"+p.Title+"* Anonymous poll created by @" + p.Owner.UserName
	} else {
		txt = "*"+p.Title+"* Poll created by @" + p.Owner.UserName
	}

	blocks = append(blocks,
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", txt, false, false),
			nil,
			nil,
		),
		slack.NewDividerBlock(),
	)

	for _, item := range p.Options {
		be := slack.NewButtonBlockElement(p.ID+"::"+item.Title, item.Title, slack.NewTextBlockObject("plain_text", "vote", false, false))
		bk := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", item.Title, false, false), nil, slack.NewAccessory(be))

		blocks = append(blocks, bk)

		var voters []slack.MixedElement
		for _, v := range item.Votes {
			if p.Mode.Anonymous {
				voters = append(voters, *slack.NewTextBlockObject("plain_text", ":thumbsup:", true, false))
			} else {
				voters = append(voters, *slack.NewImageBlockElement(v.Avatar, v.UserName))
			}
		}

		v := len(voters)
		var m string
		switch v {
		case 0:
			m = "No votes"
		case 1:
			m = "1 vote"
		default:
			m = fmt.Sprintf("%d votes", len(voters))
		}

		voters = append(voters, slack.NewTextBlockObject("plain_text", m, false, false))

		if len(voters) != 0 {
			ct := slack.NewContextBlock(item.Title+" voters", voters...)
			blocks = append(blocks, ct)
		}
	}

	blocks = append(blocks, slack.NewDividerBlock())

	msg := slack.NewBlockMessage(blocks...)
	msg.Msg.ReplaceOriginal = true
	msg.Msg.ResponseType = slack.ResponseTypeInChannel

	return &msg
}
