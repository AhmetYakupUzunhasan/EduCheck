package models

type Assignment struct {
	ID          int
	Title       string
	Explanation string
	CreatedAt   string
	ExpiresAt   string
	Status      string
}

type UserToAssignment struct {
	ID           int
	AssignmentID int
	UserID       int
	CreatedAt    string
	Status       string
}

type PostAssignment struct {
	Title       string `json:"title" binding:"required"`
	Explanation string `json:"explanation" binding:"required"`
	CreatedAt   string
	ExpiresAt   string `json:"expires_at" binding:"required"`
}
