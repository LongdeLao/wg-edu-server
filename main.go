package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"wg-edu-server/config"
	"wg-edu-server/handlers"
	"wg-edu-server/models"
	"wg-edu-server/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Load configuration
	config := config.NewConfig()

	// Connect to the database
	db, err := models.NewDB(
		config.DBHost,
		config.DBPort,
		config.DBName,
		config.DBUser,
		config.DBPassword,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database")

	// Initialize database schema
	if err := initializeDatabase(db); err != nil {
		log.Printf("Warning: failed to initialize database schema: %v", err)
	}

	// Create test users
	if err := CreateTestUsers(db); err != nil {
		log.Printf("Warning: failed to create test users: %v", err)
	}

	// Create handler with dependencies
	handler := &handlers.Handler{
		DB:        db,
		JWTSecret: config.JWTSecret,
	}

	// Setup Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, handler)

	// Start the server
	log.Printf("Server starting on port %s", config.ServerPort)
	if err := router.Run(config.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// InitializeDatabase sets up the database schema
func initializeDatabase(db *models.DB) error {
	// Read schema files
	schemaFiles := []string{
		"schema_teachers.sql",
	}

	for _, file := range schemaFiles {
		// Check if file exists
		if _, err := os.Stat(file); os.IsNotExist(err) {
			log.Printf("Schema file %s not found, skipping", file)
			continue
		}

		// Read file content
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read schema file %s: %v", file, err)
		}

		// Execute SQL
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute schema file %s: %v", file, err)
		}

		log.Printf("Successfully applied schema from %s", file)
	}

	return nil
}

// CreateTestUsers creates test users if they don't exist
func CreateTestUsers(db *models.DB) error {
	// Create generic test users
	roles := []string{"admin", "teacher", "student"}

	for _, role := range roles {
		username := fmt.Sprintf("test_%s", role)
		password := fmt.Sprintf("test_%s", role)

		// Check if user exists
		_, err := db.GetUserByUsername(username)
		if err == nil {
			// User already exists
			continue
		}

		// Create user
		_, err = db.CreateUser(username, password, role)
		if err != nil {
			return fmt.Errorf("failed to create test user '%s': %v", username, err)
		}

		log.Printf("Created test user: %s with role: %s", username, role)
	}

	// Create specific teacher users
	teachers := []string{"wg", "liz", "eddie", "yu", "tan", "li"}
	for _, username := range teachers {
		password := fmt.Sprintf("%s_123", username)

		// Check if teacher exists
		_, err := db.GetUserByUsername(username)
		if err == nil {
			// Teacher already exists
			continue
		}

		// Create teacher
		_, err = db.CreateUser(username, password, "teacher")
		if err != nil {
			return fmt.Errorf("failed to create teacher '%s': %v", username, err)
		}

		log.Printf("Created teacher: %s", username)
	}

	return nil
}
