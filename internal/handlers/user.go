package handlers

import (
	"EduCheck/internal/database"
	"EduCheck/internal/middleware"
	"EduCheck/internal/models"
	"net/http"
	"net/mail"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

const key string = "osul xjdm xjqm iqwo"

func PostUser(ctx *gin.Context) {
	var user models.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	lenUsername := utf8.RuneCountInString(user.Username)
	if lenUsername < 4 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data":  nil,
			"error": "Username can't be shorter than 4 characters",
		})
		return
	}

	lenPassword := utf8.RuneCountInString(user.Password)
	if lenPassword < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data":  nil,
			"error": "Password can't be shorter than 8 characters",
		})
		return
	}

	_, err = mail.ParseAddress(user.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	user.Password, err = middleware.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	t := time.Now().Round(0)
	currentTime := t.Format("2006-01-02 15:04")
	user.CreatedAt = currentTime

	strID, err := database.InsertUserIntoDb(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	code := middleware.GenerateAuthCode()
	mailer := middleware.NewGmailMailer("ahmetyakupuzunhasan@gmail.com", key)
	go mailer.Send(user.Email, "Verify Your Email", code)
	var emailVerification models.EmailVerification
	if emailVerification.UserID, err = strconv.Atoi(strID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	emailVerification.Code = code
	t = time.Now().Add(time.Minute * 3)
	emailVerification.ExpiresAt = t.Format("2006-01-02 15:04")
	err = database.InsertEmailVerificationIntoDb(&emailVerification)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data":    user,
		"message": "Verification Email Sent",
	})
}

func VerifyEmail(ctx *gin.Context) {
	var reqBody models.EmailVerificationPostRequest
	err := ctx.ShouldBindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	emailVerification, err := database.SelectEmailVerificationFromDb(reqBody.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	if emailVerification.Status == "resolved" {
		ctx.JSON(http.StatusOK, gin.H{
			"data":  "Already Verified",
			"error": nil,
		})
		return
	}

	t := time.Now().Round(0)
	currentTime := t.Format("2006-01-02 15:04")

	if currentTime >= emailVerification.ExpiresAt {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"data":  nil,
			"error": "Time Has Run Out",
		})
		return
	}

	if emailVerification.Code != reqBody.Code {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"data":  nil,
			"error": "Code is incorrect",
		})
		return
	}

	if err = database.UpdateEmailVerificationStateInDb(emailVerification.UserID, "resolved"); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	err = database.UpdateUserStateInDb(reqBody.UserID, "active")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	strID := strconv.Itoa(reqBody.UserID)
	token, err := middleware.GenerateToken(strID, "user")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  token,
		"error": nil,
	})
}

func Login(ctx *gin.Context) {
	var loginRequest models.LoginRequest
	err := ctx.ShouldBindJSON(&loginRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}
	lenUsername := utf8.RuneCountInString(loginRequest.Username)
	if lenUsername < 4 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data":  nil,
			"error": "Username Length Can't Be Less Than 4",
		})
		return
	}

	lenPassword := utf8.RuneCountInString(loginRequest.Password)
	if lenPassword < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data":  nil,
			"error": "Password Length Can't Be Less Than 8",
		})
		return
	}

	hashedPassword, id, role, err := database.SelectUserPasswordFromDbByUsername(loginRequest.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"data":  nil,
			"error": "Username or Password is Incorrect",
		})
		return
	}

	isSame := middleware.CompareHashedPassword(loginRequest.Password, hashedPassword)
	if !isSame {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"data":  nil,
			"error": "Username or Password is Incorrect",
		})
		return
	}

	token, err := middleware.GenerateToken(id, role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  token,
		"error": "You're Good To Go",
	})
}

func GetUsers(ctx *gin.Context) {
	users, err := database.SelectUsersFromDb()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  users,
		"error": "You're Good To Go",
	})
}
