package routes

import (
	"github.com/BLimmie/intouch-health-capstone-2019/app"
	db "github.com/BLimmie/intouch-health-capstone-2019/db/dbcode"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client = nil
var ic *db.IntouchClient = nil

func init() {
	client = db.OpenConnection()
	ic = db.CreateIntouchClient("test", client)
}

var registry = app.LoginHandler{CurrentTokens: map[string]primitive.ObjectID{}}

func Routes() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	// Get Object Data
	router.GET("/patient/:user", getPatient)
	router.GET("/provider/:user", getProvider)
	router.GET("/session/:user", getSession)
	// Submit New Data
	router.POST("/patient", addPatient)
	router.POST("/provider", addProvider)
	router.POST("/session", addSession)

	// Get Session Metrics
	router.GET("/metrics/:id", getSessionMetrics)
	router.GET("/metrics/:id/aggregate", getSessionMetricsAggregate)
	// Get Sentiment
	router.POST("/sentiment/frame", getSentimentFrame)
	router.POST("/sentiment/frame/:id", submitSentimentFrame)
	router.POST("/sentiment/text", getSentimentText)
	router.POST("/sentiment/text/:id", submitSentimentText)
	// Authenticate
	router.POST("/login", login)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}
