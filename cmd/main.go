package main

import (
	"EduCheck/internal/database"
	"EduCheck/internal/handlers"
	"EduCheck/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	if err := database.ConnectToDb(); err != nil {
		panic("Could Not Connect To Db")
	}

	if err := database.InitializeDatatables(); err != nil {
		panic("Could Not Init DTs")
	}

	app.POST("/register", handlers.PostUser)
	app.POST("/verify-email", handlers.VerifyEmail)
	app.POST("/login", handlers.Login)

	api := app.Group("/api")

	protected := api.Group("", middleware.AuthMiddleware())
	{
		protected.GET("/users", handlers.GetUsers)
		protected.GET("/assignments", handlers.GetAssignments)
		protected.GET("/associated-assignments", handlers.GetAssociatedAssignments)

		extraProtected := protected.Group("", middleware.RequireRole("teacher"))
		{
			extraProtected.POST("/assignments", handlers.PostAssignment)
		}
	}

	app.Run(":8080")
}
