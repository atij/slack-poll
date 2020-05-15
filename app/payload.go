package app

import (
	"github.com/atij/slack-poll/model"
	"github.com/gin-gonic/gin"
)

// CommandPayload struct
type commandPayload struct {
	TeamID      string `form:"team_id" json:"team_id"`
	TeamDomain  string `form:"team_domain" json:"team_domain"`
	ChannelID   string `form:"channel_id" json:"channel_id"`
	ChannelName string `form:"channel_name" json:"channel_name"`
	UserID      string `form:"user_id" json:"user_id"`
	UserName    string `form:"user_name" json:"user_name"`
	Command     string `form:"command" json:"command"`
	Text        string `form:"text" json:"text"`
	ResponseURL string `form:"response_url" json:"response_url"`
}

type payload struct {
	Data string `form:"payload"`
}

type actionPayload struct {
	Type string `json:"type"`
	User struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		TeamID   string `json:"team_id"`
	} `json:"user"`
	APIAppID  string `json:"api_app_id"`
	Token     string `json:"token"`
	Container struct {
		Type        string `json:"type"`
		MessageTs   string `json:"message_ts"`
		ChannelID   string `json:"channel_id"`
		IsEphemeral bool   `json:"is_ephemeral"`
	} `json:"container"`
	TriggerID string `json:"trigger_id"`
	Team      struct {
		ID     string `json:"id"`
		Domain string `json:"domain"`
	} `json:"team"`
	Channel struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
	ResponseURL string `json:"response_url"`
	Actions     []struct {
		ActionID string `json:"action_id"`
		BlockID  string `json:"block_id"`
		Text     struct {
			Type  string `json:"type"`
			Text  string `json:"text"`
			Emoji bool   `json:"emoji"`
		} `json:"text"`
		Value    string `json:"value"`
		Type     string `json:"type"`
		ActionTs string `json:"action_ts"`
	} `json:"actions"`
}

// Help response structure
func getHelpReponse(c *gin.Context) {
	c.JSON(200, gin.H{
		"response_type": "ephemeral",
		"blocks": []gin.H{
			{
				"type": "section",
				"text": gin.H{
					"type": "mrkdwn",
					"text": "*Adore Me poll for slack*",
				},
			},
			{
				"type": "divider",
			},
			{
				"type": "section",
				"text": gin.H{
					"type": "mrkdwn",
					"text": "*Simple poll*",
				},
			},
			{
				"type": "section",
				"text": gin.H{
					"type": "mrkdwn",
					"text": "\n/am-poll \"What's your favourite color ?\" \"Red\" \"Green\" \"Blue\" \"Yellow\"\n",
				},
			},
		},
	})
}

// Poll response ...
func getPollResponse(p *model.Poll) gin.H {

	res := gin.H{
		"response_type": "ephemeral",
		"replace_original": "true",
	}

	var blocks []gin.H

	blocks = append(blocks,
		gin.H{
			"type": "section",
			"text": gin.H{
				"type": "mrkdwn",
				"text": p.Title,
			},
		},
		gin.H{"type": "divider"},
	)

	for _, item := range p.Options {
		blocks = append(blocks, getOption(&item, p)...)
		b := getContext(item)
		if b != nil {
			blocks = append(blocks, b)
		}
	}

	blocks = append(blocks,
		gin.H{"type": "divider"},
		gin.H{
			"type": "actions",
			"elements": []gin.H{
				{
					"type": "button",
					"text": gin.H{
						"type":  "plain_text",
						"emoji": true,
						"text":  "Add a suggestion",
					},
					"value": "click_me_123",
				},
			},
		},
	)

	res["blocks"] = blocks

	return res
}

func getOption(o *model.PollOption, p *model.Poll) []gin.H {
	return []gin.H{
		{
			"type": "section",
			"text": gin.H{
				"type": "mrkdwn",
				"text": o.Title,
			},
			"accessory": gin.H{
				"type": "button",
				"text": gin.H{
					"type":  "plain_text",
					"emoji": true,
					"text":  "Vote",
				},
				"value":     o.Title,
				"action_id": p.ID + "::" + o.Title,
			},
		},
	}
}

func getContext(o model.PollOption) gin.H {
	var elements []gin.H

	for _, item := range o.Votes {
		elements = append(elements, gin.H{
			"type":      "image",
			"image_url": "https://api.slack.com/img/blocks/bkb_template_images/profile_1.png",
			"alt_text":  item.UserName,
		})
	}

	if len(elements) == 0 {
		return nil
	}

	return gin.H{
		"type":     "context",
		"elements": elements,
	}
}
