package routes

import (
	"net/http"
	"strings"

	"github.com/BLimmie/intouch-health-capstone-2019/app"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client = nil
var ic *app.IntouchClient = nil
var DBWorkers *app.WorkerHandler = nil
var OFWorkers *app.WorkerHandler = nil
var GCPWorkers *app.WorkerHandler = nil

func init() {
	client = app.OpenConnection()
	ic = app.CreateIntouchClient("intouch", client)
}

var registry = app.NewLoginHandler()

func Routes() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	DBWorkers = app.NewWorkerHandler(4)
	OFWorkers = app.NewWorkerHandler(4)
	GCPWorkers = app.NewWorkerHandler(4)
	router := gin.Default()
	// Get Object Data
	router.POST("/patient/:user", getPatient)
	router.POST("/provider/:user", getProvider)
	router.POST("/session/:id", getSession)
	router.POST("/session/:id/updatetext", getLatestTextMetrics)
	router.POST("/sessions/:userid", getSessions)
	router.POST("/associatedsessions/:proid/:patun", getAssociatedSessions)
	router.POST("/latestsession", getLatestSession)
	router.POST("/user/:token", getUserFromToken)
	// Submit New Data
	router.POST("/patient", addPatient)
	router.POST("/provider", addProvider)
	router.POST("/session", addSession)
	router.POST("/associateUser", associateUser)

	// Get Session Metrics
	router.POST("/metrics/:id", getSessionMetrics)
	router.POST("/metrics/:id/aggregate", getSessionMetricsAggregate)
	// Get Sentiment
	router.POST("/sentiment/frame", getSentimentFrame)
	router.POST("/sentiment/frame/:id", submitSentimentFrame)
	router.POST("/sentiment/text", getSentimentText)
	router.POST("/sentiment/text/:id", submitSentimentText)
	// Get Twilio Token
	router.POST("/twilio/getToken", getToken)
	// Authenticate
	router.POST("/login", login)

	router.LoadHTMLFiles("routes/js/build/index.html")
	router.GET("/*filename", func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/client") {
			c.HTML(http.StatusOK, "index.html", nil)
		} else {
			var relaPathFromPwdBuilder strings.Builder
			relaPathFromPwdBuilder.WriteString("routes/js/build")
			relaPathFromPwdBuilder.WriteString(c.Request.URL.Path)
			c.File(relaPathFromPwdBuilder.String())
		}
	})

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}
