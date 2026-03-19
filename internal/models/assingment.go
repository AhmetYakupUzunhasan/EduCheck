package models

type Class struct {
	ID          int
	Name        string
	Description string
}

type UserToClass struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id" binding:"required"`
	ClassID   int    `json:"class_id" binding:"required"`
	CreatedAt string `json:"created_at" binding:"required"`
}

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
