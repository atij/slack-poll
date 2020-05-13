package app

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

type commandPayload struct {
	Token       string `json:"token"`
	TeamID      string `json:"team_id"`
	TeamDomain  string `json:"team_domain"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_nam"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Command     string `json:"command"`
	Text        string `json:"text"`
	ResponseURL string `json:"response_url"`
}

/*
	"token"=>"XXXXX",
  	"team_id"=>"YYYY",
  	"team_domain"=>"ZZZZ",
  	"channel_id"=>"UUUU",
  	"channel_name"=>"directmessage",
  	"user_id"=>"U1234567",
  	"user_name"=>"anderson",
  	"command"=>"/congratulate",
  	"text"=>"@john for his new product release! It's brilliant!",
	  "response_url"=>"https://hooks.slack.com/commands/YYYY/DDDDD/HASH"
*/

func command(c *gin.Context) {

	// unmarshal
	var p commandPayload
	if c.BindJSON(&p) != nil {
		c.String(200, "invalid payload")
		return
	}

	// validate token
	if !isTokenValid(&p) {
		c.String(200, "Looks like your token is invalid, check application config.")
		return
	}

	// help section
	if p.Text == "help" {
		c.String(200, "Use `/poll` to create simple poll. Example: /poll \"Poll question?\" \"Option 1\" \"Option 2\" \"Option 3\"")
		return
	}

	c.JSON(200, gin.H{
		"message": "command route called!",
	})
}

func isTokenValid(p *commandPayload) bool {
	if p.Token != os.Getenv("SLACK_TOKEN") {
		return false
	}
	return true
}

// StartApplication ...
func StartApplication() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	var router = gin.New()

	router.GET("/readyz", readyz)
	router.GET("/healthz", healthz)
	router.POST("/slack/command", command)

	router.Run(":" + port)
}
