-- Create subjects table
CREATE TABLE IF NOT EXISTS subjects (
    id SERIAL PRIMARY KEY,
    grade VARCHAR(5) NOT NULL CHECK (grade IN ('PIB', 'IB1', 'IB2')),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create teacher_subjects mapping table
CREATE TABLE IF NOT EXISTS teacher_subjects (
    id SERIAL PRIMARY KEY,
    teacher_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subject_id INTEGER NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(teacher_id, subject_id)
);

-- Create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_teacher_subjects_teacher_id ON teacher_subjects(teacher_id);
CREATE INDEX IF NOT EXISTS idx_teacher_subjects_subject_id ON teacher_subjects(subject_id);
CREATE INDEX IF NOT EXISTS idx_subjects_grade ON subjects(grade);

-- Insert subjects
DO $$
BEGIN
    -- PIB Subjects
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'PIB' AND name = 'Mathematics') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('PIB', 'Mathematics', 'Pre-IB Mathematics');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'PIB' AND name = 'Additional Mathematics') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('PIB', 'Additional Mathematics', 'Pre-IB Additional Mathematics');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'PIB' AND name = 'Physics') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('PIB', 'Physics', 'Pre-IB Physics');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'PIB' AND name = 'Chemistry') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('PIB', 'Chemistry', 'Pre-IB Chemistry');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'PIB' AND name = 'Biology') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('PIB', 'Biology', 'Pre-IB Biology');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'PIB' AND name = 'English') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('PIB', 'English', 'Pre-IB English');
    END IF;
    
    -- IB1 Subjects
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'IB1' AND name = 'Math AA') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('IB1', 'Math AA', 'IB1 Mathematics Analysis and Approaches');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'IB1' AND name = 'Math AI') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('IB1', 'Math AI', 'IB1 Mathematics Applications and Interpretation');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'IB1' AND name = 'Physics') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('IB1', 'Physics', 'IB1 Physics');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'IB1' AND name = 'Chemistry') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('IB1', 'Chemistry', 'IB1 Chemistry');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'IB1' AND name = 'Business Management') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('IB1', 'Business Management', 'IB1 Business Management');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'IB1' AND name = 'Biology') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('IB1', 'Biology', 'IB1 Biology');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'IB1' AND name = 'English B') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('IB1', 'English B', 'IB1 English B');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'IB1' AND name = 'Economics') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('IB1', 'Economics', 'IB1 Economics');
    END IF;
    
    -- IB2 Subjects
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'IB2' AND name = 'Math AA') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('IB2', 'Math AA', 'IB2 Mathematics Analysis and Approaches');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'IB2' AND name = 'Math AI') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('IB2', 'Math AI', 'IB2 Mathematics Applications and Interpretation');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'IB2' AND name = 'Physics') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('IB2', 'Physics', 'IB2 Physics');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'IB2' AND name = 'Chemistry') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('IB2', 'Chemistry', 'IB2 Chemistry');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'IB2' AND name = 'Business Management') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('IB2', 'Business Management', 'IB2 Business Management');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'IB2' AND name = 'Biology') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('IB2', 'Biology', 'IB2 Biology');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'IB2' AND name = 'English B') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('IB2', 'English B', 'IB2 English B');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM subjects WHERE grade = 'IB2' AND name = 'Economics') THEN
        INSERT INTO subjects (grade, name, description) VALUES ('IB2', 'Economics', 'IB2 Economics');
    END IF;
END
$$;

-- Assign subjects to teachers
DO $$
DECLARE
    wg_id INTEGER;
    yu_id INTEGER;
    eddie_id INTEGER;
    li_id INTEGER;
    tan_id INTEGER;
    liz_id INTEGER;
    
    pib_physics_id INTEGER;
    ib1_physics_id INTEGER;
    ib2_physics_id INTEGER;
    pib_math_id INTEGER;
    ib1_math_aa_id INTEGER;
    ib1_math_ai_id INTEGER;
    ib2_math_aa_id INTEGER;
    ib2_math_ai_id INTEGER;
    ib1_economics_id INTEGER;
    ib2_economics_id INTEGER;
    pib_english_id INTEGER;
    ib1_english_id INTEGER;
    ib2_english_id INTEGER;
    pib_biology_id INTEGER;
    pib_chemistry_id INTEGER;
    ib1_business_id INTEGER;
    ib2_business_id INTEGER;
BEGIN
    -- Get teacher IDs
    SELECT id INTO wg_id FROM users WHERE username = 'wg' AND role = 'teacher';
    SELECT id INTO yu_id FROM users WHERE username = 'yu' AND role = 'teacher';
    SELECT id INTO eddie_id FROM users WHERE username = 'eddie' AND role = 'teacher';
    SELECT id INTO li_id FROM users WHERE username = 'li' AND role = 'teacher';
    SELECT id INTO tan_id FROM users WHERE username = 'tan' AND role = 'teacher';
    SELECT id INTO liz_id FROM users WHERE username = 'liz' AND role = 'teacher';
    
    -- Get subject IDs
    SELECT id INTO pib_physics_id FROM subjects WHERE grade = 'PIB' AND name = 'Physics';
    SELECT id INTO ib1_physics_id FROM subjects WHERE grade = 'IB1' AND name = 'Physics';
    SELECT id INTO ib2_physics_id FROM subjects WHERE grade = 'IB2' AND name = 'Physics';
    SELECT id INTO pib_math_id FROM subjects WHERE grade = 'PIB' AND name = 'Mathematics';
    SELECT id INTO ib1_math_aa_id FROM subjects WHERE grade = 'IB1' AND name = 'Math AA';
    SELECT id INTO ib1_math_ai_id FROM subjects WHERE grade = 'IB1' AND name = 'Math AI';
    SELECT id INTO ib2_math_aa_id FROM subjects WHERE grade = 'IB2' AND name = 'Math AA';
    SELECT id INTO ib2_math_ai_id FROM subjects WHERE grade = 'IB2' AND name = 'Math AI';
    SELECT id INTO ib1_economics_id FROM subjects WHERE grade = 'IB1' AND name = 'Economics';
    SELECT id INTO ib2_economics_id FROM subjects WHERE grade = 'IB2' AND name = 'Economics';
    SELECT id INTO pib_english_id FROM subjects WHERE grade = 'PIB' AND name = 'English';
    SELECT id INTO ib1_english_id FROM subjects WHERE grade = 'IB1' AND name = 'English B';
    SELECT id INTO ib2_english_id FROM subjects WHERE grade = 'IB2' AND name = 'English B';
    SELECT id INTO pib_biology_id FROM subjects WHERE grade = 'PIB' AND name = 'Biology';
    SELECT id INTO pib_chemistry_id FROM subjects WHERE grade = 'PIB' AND name = 'Chemistry';
    SELECT id INTO ib1_business_id FROM subjects WHERE grade = 'IB1' AND name = 'Business Management';
    SELECT id INTO ib2_business_id FROM subjects WHERE grade = 'IB2' AND name = 'Business Management';
    
    -- Assign subjects to teachers if they exist
    
    -- WG: PIB Physics, IB1 Physics, IB2 Physics, PIB Math, IB1 Math AA, IB1 Math AI, IB2 Math AA, IB2 Math AI
    IF wg_id IS NOT NULL AND pib_physics_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (wg_id, pib_physics_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
    IF wg_id IS NOT NULL AND ib1_physics_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (wg_id, ib1_physics_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
    IF wg_id IS NOT NULL AND ib2_physics_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (wg_id, ib2_physics_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
    IF wg_id IS NOT NULL AND pib_math_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (wg_id, pib_math_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
    IF wg_id IS NOT NULL AND ib1_math_aa_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (wg_id, ib1_math_aa_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
    IF wg_id IS NOT NULL AND ib1_math_ai_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (wg_id, ib1_math_ai_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
    IF wg_id IS NOT NULL AND ib2_math_aa_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (wg_id, ib2_math_aa_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
    IF wg_id IS NOT NULL AND ib2_math_ai_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (wg_id, ib2_math_ai_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
    -- YU: IB1 Economics, IB2 Economics
    IF yu_id IS NOT NULL AND ib1_economics_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (yu_id, ib1_economics_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
    IF yu_id IS NOT NULL AND ib2_economics_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (yu_id, ib2_economics_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
    -- Eddie: PIB English, IB1 English B, IB2 English B
    IF eddie_id IS NOT NULL AND pib_english_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (eddie_id, pib_english_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
    IF eddie_id IS NOT NULL AND ib1_english_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (eddie_id, ib1_english_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
    IF eddie_id IS NOT NULL AND ib2_english_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (eddie_id, ib2_english_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
    -- Li: PIB Biology
    IF li_id IS NOT NULL AND pib_biology_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (li_id, pib_biology_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
    -- Tan: PIB Chemistry
    IF tan_id IS NOT NULL AND pib_chemistry_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (tan_id, pib_chemistry_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
    -- Liz: IB1 Business Management, IB2 Business Management
    IF liz_id IS NOT NULL AND ib1_business_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (liz_id, ib1_business_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
    IF liz_id IS NOT NULL AND ib2_business_id IS NOT NULL THEN
        INSERT INTO teacher_subjects (teacher_id, subject_id) 
        VALUES (liz_id, ib2_business_id)
        ON CONFLICT (teacher_id, subject_id) DO NOTHING;
    END IF;
    
END
$$; 