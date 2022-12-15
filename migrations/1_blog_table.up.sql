CREATE TABLE IF NOT EXISTS "users"(
    "id" SERIAL PRIMARY KEY,
    "first_name" VARCHAR(30) NOT NULL,
    "last_name" VARCHAR(30) NOT NULL,
    "phone_number" VARCHAR(20) UNIQUE,
    "email" VARCHAR(50) NOT NULL UNIQUE,
    "gender" VARCHAR(10) CHECK ("gender" IN('male', 'female')),
    "password" VARCHAR NOT NULL,
    "username" VARCHAR(30) UNIQUE,
    "profile_image_url" VARCHAR,
    "type" VARCHAR(255) CHECK ("type" IN('superadmin', 'user')) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Password; zohid2004
INSERT INTO users(first_name, last_name, email, password, type) 
VALUES('zohid', 'saidov', 'zohidsaidov17@gmail.com', '$2a$10$nbFlsRUTs.8//V3s5TAczuETJClXKWx9zM95Kami9n9V1ODL/tQti', 'superadmin');