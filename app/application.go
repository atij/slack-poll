package app

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/atij/slack-poll/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func initDB() (*repository.FirestoreRepository, error) {
	c, err := firestore.NewClient(context.Background(), os.Getenv("GOOGLE_APPLICATION_PROJECT_ID"))
	if err != nil {
		return nil, err
	}

	db, err := repository.NewFirestoreRepository(c, os.Getenv("FIRESTORE_COLLECTION"))
	return db, nil
}

// StartApplication ...
func StartApplication() {

	// TODO: move config to separate function
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	db, err := initDB()
	if err != nil {
		log.Fatalf("cannot init firestore database, error: %v", err)
	}

	var router = gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	/*
		router.Use(func(c *gin.Context) {
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ := ioutil.ReadAll(tee)
			c.Request.Body = ioutil.NopCloser(&buf)
			log.Print(string(body))
			log.Print(c.Request.Header)
			c.Next()
		})
	*/

	router.GET("/readyz", readyz)
	router.GET("/healthz", healthz)
	router.POST("/slack/command", command)
	router.POST("/slack/poll/actions", pollActions)

	router.Run(":" + port)
}
