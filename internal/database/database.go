package database

import (
	"EduCheck/internal/models"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB

func ConnectToDb() error {
	var err error
	db, err = sql.Open("sqlite", "./app.db")
	if err != nil {
		log.Fatal("Error when connecting to db: ", err.Error())
	}

	return nil
}

func InitializeDatatables() error {
	//Users
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password NOT NULL,
			created_at TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'active')),
			role TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin'))
		);
	`)
	if err != nil {
		log.Fatal("Error when initializing users datatable")
	}

	//Email Verification
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS email_verification(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL UNIQUE,
			code TEXT NOT NULL UNIQUE,
			expires_at TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'unresolved' CHECK (status IN ('unresolved', 'resolved')),

			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		log.Fatal("Error when initializing email_verification datatable: ", err)
	}

	//Assignment
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS assignments(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			explanation TEXT NOT NULL,
			created_at TEXT NOT NULL,
			expires_at TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive'))
		);
	`)
	if err != nil {
		log.Fatal("Error when initializing assignments datatable: ", err)
	}

	//User To Assignment
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_to_assignment(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			assignment_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			created_at TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'responded')),

			FOREIGN KEY (assignment_id) REFERENCES assignments(id) ON DELETE CASCADE
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		log.Fatal("Error when initializing user to assignment datatable: ", err)
	}

	return nil
}

func InsertAssignmentIntoDb(assignment *models.PostAssignment) error {
	sql := "INSERT INTO assignments (title, explanation, created_at, expires_at) VALUES (?,?,?,?)"
	if _, err := db.Exec(sql, assignment.Title, assignment.Explanation, assignment.CreatedAt, assignment.ExpiresAt); err != nil {
		fmt.Println("Error when inserting assignment into db: ", err)
		return err
	}

	return nil
}

func SelectUserToAssignmentFromDbByUserId(userID string) ([]int, error) {
	query := "SELECT assignment_id FROM user_to_assignment WHERE user_id = ?"
	rows, err := db.Query(query, userID)
	if err != nil {
		fmt.Println("Error when selecting user to assignment from db by user_id: ", err)
		return nil, err
	}
	defer rows.Close()

	var userToAssignments []int
	for rows.Next() {
		var c int
		if err := rows.Scan(&c); err != nil {
			fmt.Println("Error when selecting user to assignment from db by user_id: ", err)
			return nil, err
		}
		userToAssignments = append(userToAssignments, c)
	}

	return userToAssignments, nil
}

func SelectAssignmentsByIdsFromDb(ids []int) ([]models.Assignment, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("No Things There")
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := range ids {
		placeholders[i] = "?"
		args[i] = ids[i]
	}

	query := fmt.Sprintf("SELECT * FROM assignments WHERE id IN (%s)", strings.Join(placeholders, ","))

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println("Error when selecting assignment from db by id: ", err)
		return nil, err
	}
	defer rows.Close()

	var assignments []models.Assignment
	for rows.Next() {
		var c models.Assignment
		if err := rows.Scan(&c.ID, &c.Title, &c.Explanation, &c.CreatedAt, &c.ExpiresAt, &c.Status); err != nil {
			fmt.Println("Error when selecting assignment from db by id: ", err)
			return nil, err
		}
		assignments = append(assignments, c)
	}

	return assignments, nil

}

func SelectAssignmentsFromDb() ([]models.Assignment, error) {
	query := "SELECT * FROM assignments"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error when selecting assignments from db: ", err)
		return nil, err
	}
	defer rows.Close()

	var assignments []models.Assignment
	for rows.Next() {
		var c models.Assignment
		if err := rows.Scan(&c.ID, &c.Title, &c.Explanation, &c.CreatedAt, &c.ExpiresAt, &c.Status); err != nil {
			return nil, err
		}
		assignments = append(assignments, c)
	}

	return assignments, nil

}

func InsertUserIntoDb(user *models.User) (string, error) {
	query := "INSERT INTO users (username, email, password, created_at) VALUES (?,?,?,?)"
	result, err := db.Exec(query, user.Username, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		fmt.Println("Error when inserting user into db: ", err)
		return "", err
	}

	ID, err := result.LastInsertId()
	strID := strconv.Itoa(int(ID))
	return strID, nil
}

func SelectUserPasswordFromDbByUsername(username string) (string, string, string, error) {
	query := "SELECT id, password, role FROM users WHERE username = ?"
	row := db.QueryRow(query, username)
	var id, password, role string
	err := row.Scan(&id, &password, &role)
	if err != nil {
		fmt.Println("Error when selecing user password from db by username")
		return "", "", "", err
	}

	return password, id, role, nil
}

func SelectUsersFromDb() ([]models.User, error) {
	query := "SELECT * FROM users"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error when selecting users from db: ", err)
		return nil, err
	}

	var users []models.User
	for rows.Next() {
		var c models.User
		err = rows.Scan(&c.ID, &c.Username, &c.Email, &c.Password, &c.CreatedAt, &c.Status, &c.Role)
		if err != nil {
			fmt.Println("Error when scaning users from db: ", err)
			return nil, err
		}
		users = append(users, c)
	}

	return users, nil
}

func UpdateUserStateInDb(id int, state string) error {
	query := "UPDATE users SET status = ? WHERE id = ?"
	_, err := db.Exec(query, state, id)
	if err != nil {
		fmt.Println("Error when updating user state: ", err)
		return err
	}

	return nil
}

func InsertEmailVerificationIntoDb(emailVerification *models.EmailVerification) error {
	query := "INSERT INTO email_verification (user_id, code, expires_at) VALUES (?,?,?)"
	_, err := db.Exec(query, emailVerification.UserID, emailVerification.Code, emailVerification.ExpiresAt)
	if err != nil {
		fmt.Println("Error when insterting email verification into db: ", err)
		return err
	}

	return nil
}

func SelectEmailVerificationFromDb(userID int) (*models.EmailVerification, error) {
	query := "SELECT * FROM email_verification WHERE user_id = ?"
	row := db.QueryRow(query, userID)
	emailVerification := &models.EmailVerification{}
	err := row.Scan(&emailVerification.ID, &emailVerification.UserID, &emailVerification.Code, &emailVerification.ExpiresAt, &emailVerification.Status)
	if err != nil {
		fmt.Println("Error when selecting email_verification by user id from db: ", err)
		return nil, err
	}

	return emailVerification, nil
}

func UpdateEmailVerificationStateInDb(userID int, state string) error {
	query := "UPDATE email_verification SET status = ? WHERE user_id = ?"
	_, err := db.Exec(query, state, userID)
	if err != nil {
		fmt.Println("Error when updating email verification status in db: ", err)
		return err
	}

	return nil
}
