# WG Education Server

The backend API server for the WG Education platform, built with Go and Gin framework.

## Architecture

The application follows a clean architecture approach with the following components:

- **Handlers**: HTTP request handlers that process incoming API requests
- **Models**: Data structures and database operations
- **Routes**: API endpoint definitions and middleware
- **Config**: Application configuration management

## API Endpoints

### Authentication
- `POST /api/login` - Authenticate user and get JWT token
- `GET /api/health` - Health check endpoint

### Student Management (Admin only)
- `GET /api/admin/students` - Get all students
- `GET /api/admin/students/:id` - Get a specific student
- `POST /api/admin/students` - Create a new student
- `PUT /api/admin/students/:id` - Update a student
- `DELETE /api/admin/students/:id` - Delete a student

## Database Schema

The application uses PostgreSQL with the following schema:

### Users Table
Stores authentication information for all users:
- `id`: Serial primary key
- `username`: Unique username
- `password`: User password (currently stored as plaintext)
- `role`: User role (admin, teacher, student)
- `date_created`: Timestamp of user creation

### Students Table
Stores additional information for student users:
- `id`: Serial primary key
- `user_id`: Foreign key to users table
- `first_name`: Student's first name
- `last_name`: Student's last name
- `email`: Student's email address
- `grade`: Student's grade/class
- `created_at`: Timestamp of record creation
- `updated_at`: Timestamp of last update

## Setup and Installation

### Prerequisites
- Go 1.16+
- PostgreSQL database

### Installation

1. Clone the repository:
```
git clone https://github.com/username/wg-edu.git
cd wg-edu/wg-edu-server
```

2. Install dependencies:
```
go mod tidy
```

3. Configure database:
```
psql -h YOUR_DB_HOST -U YOUR_DB_USER -d YOUR_DB_NAME -f schema.sql
psql -h YOUR_DB_HOST -U YOUR_DB_USER -d YOUR_DB_NAME -f schema_students.sql
```

4. Run the server:
```
go run main.go
```

## Development

### Adding New Features

1. Update models if needed (models package)
2. Create new handler functions (handlers package)
3. Add routes for new endpoints (routes package)
4. Update main.go if configuration changes are needed

### Security Notes

1. Current implementation uses plaintext passwords for simplicity. In production, always use password hashing.
2. JWT tokens should have appropriate expiration times and be securely stored.

## Testing

Run tests with:
```
go test ./...
```

## Deployment

The application can be compiled into a single binary for easy deployment:

```
go build -o wg-edu-server
``` 