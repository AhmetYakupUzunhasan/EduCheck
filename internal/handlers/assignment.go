package handlers

import (
	"EduCheck/internal/database"
	"EduCheck/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func PostAssignment(ctx *gin.Context) {
	var assignment models.PostAssignment
	if err := ctx.ShouldBindJSON(&assignment); err != nil {
		ctx.JSON(400, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	expirationDate, err := strconv.Atoi(assignment.ExpiresAt)
	if err != nil {
		ctx.JSON(400, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	t := time.Now().Round(0)
	assignment.CreatedAt = t.Format("2006-01-02 15:04")

	t2 := time.Now().Add(time.Duration(expirationDate) * 24 * time.Hour).Round(0)
	assignment.ExpiresAt = t2.Format("2006-01-02 15:04")

	if err := database.InsertAssignmentIntoDb(&assignment); err != nil {
		ctx.JSON(500, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{
		"data":  assignment,
		"error": nil,
	})
}
