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

	// unmarshal
	var p commandPayload
	if c.Bind(&p) != nil {
		c.String(200, "invalid payload")
		return
	}

	// help section
	if p.Text == "help" {
		getHelpReponse(c)
		return
	}

	// clear quoutes
	p.Text = cleanDoubleQuotes(p.Text)

	// split options
	options := splitOptions(p.Text)

	// create poll
	poll, err := createPoll(&p, options)
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

	var cp actionPayload
	json.Unmarshal([]byte(p.Data), &cp)

	// get poll id
	if len(cp.Actions) == 0 {
		c.JSON(400, gin.H{
			"message": "empty actions  not allowed",
		})
		return
	}

	a := cp.Actions[0]
	pollID := strings.Split(a.ActionID, "::")[0]

	db := c.MustGet("db").(*repository.FirestoreRepository)
	poll, err := db.Find(pollID)
	if err != nil {
		c.JSON(401, gin.H{
			"message": "poll not found",
		})
		return
	}

	v := model.Vote{
		UserID:   cp.User.ID,
		UserName: cp.User.Username,
	}

	poll.AddVote(a.Value, v)

	err = db.Update(poll.ID, poll)
	if err != nil {
		c.JSON(404, err)
		return
	}

	res, _ := json.Marshal(getPollResponse(poll))
	http.Post(cp.ResponseURL, "application/json", bytes.NewBuffer(res))

	c.JSON(200, gin.H{
		"messaage": "thank you",
	})
}
