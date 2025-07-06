-- Create students table
CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    grade VARCHAR(20),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_students_user_id ON students(user_id);
CREATE INDEX IF NOT EXISTS idx_students_email ON students(email);

-- Insert some sample students (only if they don't exist)
DO $$
BEGIN
    -- First check if we have the test_student user
    DECLARE student_user_id INTEGER;
    BEGIN
        SELECT id INTO student_user_id FROM users WHERE username = 'test_student';
        
        -- If we found the test_student user but no record in students table
        IF student_user_id IS NOT NULL AND 
           NOT EXISTS (SELECT 1 FROM students WHERE user_id = student_user_id) THEN
            
            INSERT INTO students (user_id, first_name, last_name, email, grade)
            VALUES (student_user_id, 'Test', 'Student', 'test.student@example.com', 'IB1');
            
        END IF;
    EXCEPTION
        WHEN NO_DATA_FOUND THEN
            -- No test_student user, do nothing
            NULL;
    END;
END
$$; 