// Package models provides database models and operations for the WG Education platform.
package models

import (
	"database/sql"
	"time"
)

// Subject represents a subject that can be taught by teachers
type Subject struct {
	ID          int       `json:"id"`          // Unique identifier
	Grade       string    `json:"grade"`       // Educational grade (PIB, IB1, or IB2)
	Name        string    `json:"name"`        // Subject name
	Description string    `json:"description"` // Subject description
	CreatedAt   time.Time `json:"created_at"`  // Creation timestamp
}

// Teacher represents a teacher with their associated subjects
type Teacher struct {
	ID        int       `json:"id"`         // User ID from the users table
	Username  string    `json:"username"`   // Login username
	FirstName string    `json:"first_name"` // Teacher's first name
	LastName  string    `json:"last_name"`  // Teacher's last name
	Email     string    `json:"email"`      // Teacher's email address
	Subjects  []Subject `json:"subjects"`   // Subjects taught by this teacher
}

// TeacherSubject represents a mapping between teachers and subjects
type TeacherSubject struct {
	ID        int       `json:"id"`         // Unique identifier
	TeacherID int       `json:"teacher_id"` // Reference to users table
	SubjectID int       `json:"subject_id"` // Reference to subjects table
	CreatedAt time.Time `json:"created_at"` // Creation timestamp
}

// GetAllSubjects retrieves all subjects
//
// Returns:
//   - []*Subject: Array of all subjects
//   - error: Error if retrieval fails
func (db *DB) GetAllSubjects() ([]*Subject, error) {
	query := `
		SELECT id, grade, name, description, created_at
		FROM subjects
		ORDER BY grade, name
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []*Subject
	for rows.Next() {
		subject := &Subject{}
		err := rows.Scan(
			&subject.ID,
			&subject.Grade,
			&subject.Name,
			&subject.Description,
			&subject.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subjects, nil
}

// GetSubjectsByGrade retrieves all subjects for a specific grade
//
// Parameters:
//   - grade: Educational grade (PIB, IB1, or IB2)
//
// Returns:
//   - []*Subject: Array of subjects for the specified grade
//   - error: Error if retrieval fails
func (db *DB) GetSubjectsByGrade(grade string) ([]*Subject, error) {
	query := `
		SELECT id, grade, name, description, created_at
		FROM subjects
		WHERE grade = $1
		ORDER BY name
	`
	rows, err := db.Query(query, grade)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []*Subject
	for rows.Next() {
		subject := &Subject{}
		err := rows.Scan(
			&subject.ID,
			&subject.Grade,
			&subject.Name,
			&subject.Description,
			&subject.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subjects, nil
}

// GetSubjectByID retrieves a subject by ID
//
// Parameters:
//   - id: Subject ID to retrieve
//
// Returns:
//   - *Subject: Subject if found
//   - error: Error if subject not found or database error
func (db *DB) GetSubjectByID(id int) (*Subject, error) {
	subject := &Subject{}
	query := `
		SELECT id, grade, name, description, created_at
		FROM subjects
		WHERE id = $1
	`
	err := db.QueryRow(query, id).Scan(
		&subject.ID,
		&subject.Grade,
		&subject.Name,
		&subject.Description,
		&subject.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return subject, nil
}

// GetAllTeachers retrieves all teachers with their associated subjects
//
// Returns:
//   - []*Teacher: Array of all teachers with their subjects
//   - error: Error if retrieval fails
func (db *DB) GetAllTeachers() ([]*Teacher, error) {
	// Get all users with role 'teacher'
	query := `
		SELECT id, username
		FROM users
		WHERE role = 'teacher'
		ORDER BY username
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []*Teacher
	for rows.Next() {
		teacher := &Teacher{}
		err := rows.Scan(
			&teacher.ID,
			&teacher.Username,
		)
		if err != nil {
			return nil, err
		}

		// Get subjects for this teacher
		subjects, err := db.GetTeacherSubjects(teacher.ID)
		if err != nil {
			return nil, err
		}
		teacher.Subjects = subjects

		teachers = append(teachers, teacher)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return teachers, nil
}

// GetTeacherByID retrieves a teacher by ID with their associated subjects
//
// Parameters:
//   - id: Teacher ID to retrieve
//
// Returns:
//   - *Teacher: Teacher with their subjects if found
//   - error: Error if teacher not found or database error
func (db *DB) GetTeacherByID(id int) (*Teacher, error) {
	teacher := &Teacher{}
	query := `
		SELECT id, username
		FROM users
		WHERE id = $1 AND role = 'teacher'
	`
	err := db.QueryRow(query, id).Scan(
		&teacher.ID,
		&teacher.Username,
	)
	if err != nil {
		return nil, err
	}

	// Get subjects for this teacher
	subjects, err := db.GetTeacherSubjects(teacher.ID)
	if err != nil {
		return nil, err
	}
	teacher.Subjects = subjects

	return teacher, nil
}

// GetTeacherSubjects retrieves all subjects taught by a specific teacher
//
// Parameters:
//   - teacherID: Teacher ID to retrieve subjects for
//
// Returns:
//   - []Subject: Array of subjects taught by the teacher
//   - error: Error if retrieval fails
func (db *DB) GetTeacherSubjects(teacherID int) ([]Subject, error) {
	query := `
		SELECT s.id, s.grade, s.name, s.description, s.created_at
		FROM subjects s
		JOIN teacher_subjects ts ON s.id = ts.subject_id
		WHERE ts.teacher_id = $1
		ORDER BY s.grade, s.name
	`
	rows, err := db.Query(query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []Subject
	for rows.Next() {
		subject := Subject{}
		err := rows.Scan(
			&subject.ID,
			&subject.Grade,
			&subject.Name,
			&subject.Description,
			&subject.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subjects, nil
}

// AssignSubjectToTeacher assigns a subject to a teacher
//
// Parameters:
//   - teacherID: Teacher ID
//   - subjectID: Subject ID
//
// Returns:
//   - error: Error if assignment fails
func (db *DB) AssignSubjectToTeacher(teacherID, subjectID int) error {
	// First verify this is a valid teacher and subject
	var teacherExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND role = 'teacher')", teacherID).Scan(&teacherExists)
	if err != nil {
		return err
	}
	if !teacherExists {
		return sql.ErrNoRows
	}

	var subjectExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM subjects WHERE id = $1)", subjectID).Scan(&subjectExists)
	if err != nil {
		return err
	}
	if !subjectExists {
		return sql.ErrNoRows
	}

	// Create the assignment
	_, err = db.Exec(
		"INSERT INTO teacher_subjects (teacher_id, subject_id) VALUES ($1, $2) ON CONFLICT (teacher_id, subject_id) DO NOTHING",
		teacherID, subjectID,
	)
	return err
}

// RemoveSubjectFromTeacher removes a subject assignment from a teacher
//
// Parameters:
//   - teacherID: Teacher ID
//   - subjectID: Subject ID
//
// Returns:
//   - error: Error if removal fails
func (db *DB) RemoveSubjectFromTeacher(teacherID, subjectID int) error {
	_, err := db.Exec(
		"DELETE FROM teacher_subjects WHERE teacher_id = $1 AND subject_id = $2",
		teacherID, subjectID,
	)
	return err
}
