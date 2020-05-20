package app

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/atij/slack-poll/model"
	"github.com/atij/slack-poll/repository"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
)

func readyz(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ready",
	})
}

func healthz(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "healthy",
	})
}

func command(c *gin.Context) {
	sc, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		c.String(400, "failed to parse slash command")
		return
	}

	// help section
	if sc.Text == "help" {
		c.JSON(200, getHelpReponse(c))
		return
	}

	// clear quoutes
	sc.Text = cleanDoubleQuotes(sc.Text)

	// split options
	options := splitOptions(sc.Text)

	// create poll
	poll, err := createPoll(&sc, options)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err,
		})
	}

	// create poll in db
	db := c.MustGet("db").(*repository.FirestoreRepository)
	if err := db.Create(poll); err != nil {
		c.JSON(200, gin.H{
			"error": err,
		})
	}

	c.JSON(200, getPollResponse(poll))
}

func pollActions(c *gin.Context) {

	// unmarshal
	var p payload
	if err := c.Bind(&p); err != nil {
		log.Print(err)
		c.String(400, "invalid payload")
		return
	}

	p.Data, _ = url.QueryUnescape(p.Data)

	var ic slack.InteractionCallback

	json.Unmarshal([]byte(p.Data), &ic)

	if len(ic.ActionCallback.BlockActions) == 0 {
		c.JSON(400, gin.H{"message": "empty actions  not allowed"})
		return
	}

	a := ic.ActionCallback.BlockActions[0]
	pollID := strings.Split(a.ActionID, "::")[0]

	db := c.MustGet("db").(*repository.FirestoreRepository)
	poll, err := db.Find(pollID)
	if err != nil {
		c.JSON(401, gin.H{"message": "poll not found"})
		return
	}

	profile, err := getUser(ic.User.ID)
	if err != nil {
		c.JSON(403, err)
		return
	}


	v := model.Vote{
		UserID:   ic.User.ID,
		UserName: profile.RealName,
		Avatar: profile.Image24,
	}

	if poll.HasVote(a.Value, v) {
		poll.RemoveVote(a.Value, v)
	} else {
		poll.AddVote(a.Value, v)
	}

	err = db.Update(poll.ID, poll)
	if err != nil {
		c.JSON(404, err)
		return
	}

	res, _ := json.Marshal(getPollResponse(poll))
	http.Post(ic.ResponseURL, "application/json", bytes.NewBuffer(res))

	c.JSON(200, "success!")
}