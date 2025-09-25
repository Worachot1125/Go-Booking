package routes

import (
	"app/internal/logger"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Router sets up all the routes for the application
func Router(app *gin.Engine) {

	// CORS Middleware
	app.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Middleware
	app.Use(otelgin.Middleware(viper.GetString("APP_NAME")))

	// Health check endpoint
	app.GET("/healthz", func(ctx *gin.Context) {
		logger.Infof("Health check passed")
		ctx.JSON(http.StatusOK, gin.H{"status": "Health check passed.", "message": "Welcome to Project-k API."})
	})

	// Create a new group for /api/v1
	apiV1 := app.Group("/api/v1")

	// Define groups of routes under /api/v1
	Product(apiV1.Group("/products"))
	User(apiV1.Group("/users"))
	Room(apiV1.Group("/rooms"))
	Position(apiV1.Group("/positions"))
	Role(apiV1.Group("/roles"))
	Permission(apiV1.Group("/permissions"))
	Role_Permission(apiV1.Group("/role_permissions"))
	Building(apiV1.Group("/buildings"))
	Building_Room(apiV1.Group("/buildingRooms"))
	Login(apiV1.Group("/login"))
	Logout(apiV1.Group("/logout"))
	Booking(apiV1.Group("/bookings"))
	User_Role(apiV1.Group("/userRoles"))
	Equipment(apiV1.Group("/equipments"))
}
