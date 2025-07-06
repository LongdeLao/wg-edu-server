// Package models provides database models and operations for the WG Education platform.
//
// It contains database connection management, entity definitions, and data access methods.
// This package handles all interactions with the PostgreSQL database.
package models

import (
	"database/sql"
	"fmt"
	"time"
)

// User represents a user in the system with authentication information.
// Passwords are stored as plaintext for development simplicity.
// In production, passwords should be hashed using a secure algorithm.
type User struct {
	ID          int       `json:"id"`           // Unique identifier
	Username    string    `json:"username"`     // Login username
	Password    string    `json:"-"`            // Password (not included in JSON)
	Role        string    `json:"role"`         // User role: admin, teacher, or student
	DateCreated time.Time `json:"date_created"` // Account creation timestamp
}

// Student represents a student in the system with additional details.
// A student is associated with a user account for authentication.
type Student struct {
	ID        int       `json:"id"`         // Unique identifier
	UserID    int       `json:"user_id"`    // Foreign key to users table
	FirstName string    `json:"first_name"` // Student's first name
	LastName  string    `json:"last_name"`  // Student's last name
	Email     string    `json:"email"`      // Student's email address
	Grade     string    `json:"grade"`      // Student's grade/class
	CreatedAt time.Time `json:"created_at"` // Record creation timestamp
	UpdatedAt time.Time `json:"updated_at"` // Last update timestamp
	Username  string    `json:"username"`   // User's username (added for convenience)
}

// StudentRequest is used for creating or updating a student.
// Password is optional when updating.
type StudentRequest struct {
	FirstName string `json:"first_name"`         // Student's first name
	LastName  string `json:"last_name"`          // Student's last name
	Email     string `json:"email"`              // Student's email address
	Grade     string `json:"grade"`              // Student's grade/class
	Username  string `json:"username"`           // Login username
	Password  string `json:"password,omitempty"` // Login password (optional for updates)
}

// DB represents a database connection with query methods.
// It wraps the standard sql.DB connection.
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection using the provided parameters.
//
// Parameters:
//   - host: Database server hostname
//   - port: Database server port
//   - dbname: Database name
//   - user: Database username
//   - password: Database password
//
// Returns:
//   - *DB: Database connection wrapper
//   - error: Error if connection fails
func NewDB(host, port, dbname, user, password string) (*DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port, dbname, user, password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// GetUserByUsername retrieves a user by their username.
//
// Parameters:
//   - username: Username to look up
//
// Returns:
//   - *User: User object if found
//   - error: Error if user not found or database error
func (db *DB) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	query := `SELECT id, username, password, role, date_created FROM users WHERE username = $1`

	err := db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.DateCreated,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// CreateUser adds a new user to the database.
//
// Parameters:
//   - username: Login username
//   - password: Login password
//   - role: User role (admin, teacher, student)
//
// Returns:
//   - *User: Created user object
//   - error: Error if creation fails
func (db *DB) CreateUser(username, password, role string) (*User, error) {
	// Store the plain password since that's what's in the database
	user := &User{}
	query := `INSERT INTO users (username, password, role, date_created) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id, username, password, role, date_created`

	err := db.QueryRow(
		query,
		username,
		password, // Storing plaintext password
		role,
		time.Now(),
	).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.DateCreated,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// CheckPassword verifies a user's password.
//
// Parameters:
//   - password: Password to check
//
// Returns:
//   - bool: True if password matches, false otherwise
func (user *User) CheckPassword(password string) bool {
	// In this case, we're comparing plaintext passwords directly
	return user.Password == password
}

// GetAllStudents retrieves all students with their user information.
//
// Returns:
//   - []*Student: Array of all students
//   - error: Error if retrieval fails
func (db *DB) GetAllStudents() ([]*Student, error) {
	query := `
		SELECT s.id, s.user_id, s.first_name, s.last_name, s.email, s.grade, 
		       s.created_at, s.updated_at, u.username
		FROM students s
		JOIN users u ON s.user_id = u.id
		ORDER BY s.last_name, s.first_name
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*Student
	for rows.Next() {
		student := &Student{}
		err := rows.Scan(
			&student.ID,
			&student.UserID,
			&student.FirstName,
			&student.LastName,
			&student.Email,
			&student.Grade,
			&student.CreatedAt,
			&student.UpdatedAt,
			&student.Username,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// GetStudentByID retrieves a student by ID.
//
// Parameters:
//   - id: Student ID to retrieve
//
// Returns:
//   - *Student: Student object if found
//   - error: Error if student not found or database error
func (db *DB) GetStudentByID(id int) (*Student, error) {
	student := &Student{}
	query := `
		SELECT s.id, s.user_id, s.first_name, s.last_name, s.email, s.grade, 
		       s.created_at, s.updated_at, u.username
		FROM students s
		JOIN users u ON s.user_id = u.id
		WHERE s.id = $1
	`
	err := db.QueryRow(query, id).Scan(
		&student.ID,
		&student.UserID,
		&student.FirstName,
		&student.LastName,
		&student.Email,
		&student.Grade,
		&student.CreatedAt,
		&student.UpdatedAt,
		&student.Username,
	)
	if err != nil {
		return nil, err
	}
	return student, nil
}

// CreateStudent creates a new student and corresponding user.
// This operation is performed in a transaction to ensure data consistency.
//
// Parameters:
//   - req: Student creation request with personal and login information
//
// Returns:
//   - *Student: Created student object
//   - error: Error if creation fails
func (db *DB) CreateStudent(req *StudentRequest) (*Student, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Create user first
	var userID int
	userQuery := `
		INSERT INTO users (username, password, role, date_created) 
		VALUES ($1, $2, 'student', $3) 
		RETURNING id
	`
	err = tx.QueryRow(
		userQuery,
		req.Username,
		req.Password,
		time.Now(),
	).Scan(&userID)
	if err != nil {
		return nil, err
	}

	// Create student record
	now := time.Now()
	student := &Student{}
	studentQuery := `
		INSERT INTO students (user_id, first_name, last_name, email, grade, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id, user_id, first_name, last_name, email, grade, created_at, updated_at
	`
	err = tx.QueryRow(
		studentQuery,
		userID,
		req.FirstName,
		req.LastName,
		req.Email,
		req.Grade,
		now,
		now,
	).Scan(
		&student.ID,
		&student.UserID,
		&student.FirstName,
		&student.LastName,
		&student.Email,
		&student.Grade,
		&student.CreatedAt,
		&student.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	student.Username = req.Username

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return student, nil
}

// UpdateStudent updates an existing student's information.
// This operation is performed in a transaction to ensure data consistency.
//
// Parameters:
//   - id: Student ID to update
//   - req: Student update request with the new information
//
// Returns:
//   - *Student: Updated student object
//   - error: Error if update fails
func (db *DB) UpdateStudent(id int, req *StudentRequest) (*Student, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// First get the existing student to get the user_id
	var userID int
	err = tx.QueryRow("SELECT user_id FROM students WHERE id = $1", id).Scan(&userID)
	if err != nil {
		return nil, err
	}

	// Update user information if password is provided
	if req.Password != "" {
		_, err = tx.Exec("UPDATE users SET password = $1 WHERE id = $2",
			req.Password, userID)
		if err != nil {
			return nil, err
		}
	}

	// Update student information
	now := time.Now()
	student := &Student{}
	studentQuery := `
		UPDATE students 
		SET first_name = $1, last_name = $2, email = $3, grade = $4, updated_at = $5
		WHERE id = $6
		RETURNING id, user_id, first_name, last_name, email, grade, created_at, updated_at
	`
	err = tx.QueryRow(
		studentQuery,
		req.FirstName,
		req.LastName,
		req.Email,
		req.Grade,
		now,
		id,
	).Scan(
		&student.ID,
		&student.UserID,
		&student.FirstName,
		&student.LastName,
		&student.Email,
		&student.Grade,
		&student.CreatedAt,
		&student.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get username
	err = tx.QueryRow("SELECT username FROM users WHERE id = $1", userID).Scan(&student.Username)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return student, nil
}

// DeleteStudent deletes a student and their user account.
// This operation is performed in a transaction to ensure data consistency.
//
// Parameters:
//   - id: Student ID to delete
//
// Returns:
//   - error: Error if deletion fails
func (db *DB) DeleteStudent(id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get user_id first
	var userID int
	err = tx.QueryRow("SELECT user_id FROM students WHERE id = $1", id).Scan(&userID)
	if err != nil {
		return err
	}

	// Delete student record first (to maintain referential integrity)
	_, err = tx.Exec("DELETE FROM students WHERE id = $1", id)
	if err != nil {
		return err
	}

	// Delete user record
	_, err = tx.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
