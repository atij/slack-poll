package app

import (
	"os"

	"github.com/gin-gonic/gin"
)

func tokenValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.FullPath() == "/slack/command" {

			var p commandPayload

			if c.BindJSON(&p) != nil {
				c.JSON(400, gin.H{"message": "invalid json payload"})
				c.Abort()
				return
			}

			if p.Token != os.Getenv("SLACK_TOKEN") {
				c.JSON(400, gin.H{"message": "invalid token"})
				c.Abort()
				return
			}

		}
		c.Next()
	}
}
