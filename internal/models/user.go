package models

//User Needs Class Dawg
type User struct {
	ID        *int64 `json:"id,omitempty"`
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	CreatedAt string `json:"created_at,omitempty"`
	Status    string `json:"status,omitempty"`
	Role      string `json:"role,omitempty"`
}

type EmailVerification struct {
	ID        int64  `json:"id"`
	UserID    int    `json:"user_id"`
	Code      int    `json:"code"`
	ExpiresAt string `json:"expires_at"`
	Status    string `json:"status" binding:"required"`
}

type EmailVerificationPostRequest struct {
	UserID int `json:"user_id" binding:"required"`
	Code   int `json:"code" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
