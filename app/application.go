package app

import (
	"log"

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
	c.JSON(200, gin.H{
		"message": "command route called!",
	})
}

// StartApplication ...
func StartApplication() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var router = gin.Default()
	router.Use(tokenValidator())

	router.GET("/readyz", readyz)
	router.GET("/healthz", healthz)
	router.POST("/slack/command", command)

	router.Run(":8080")
}
