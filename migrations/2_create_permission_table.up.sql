CREATE TABLE IF NOT EXISTS "permissions" (
    "id" SERIAL PRIMARY KEY,
    "user_type" VARCHAR(255) CHECK ("user_type" IN('superadmin', 'user')) NOT NULL,
    "resource" VARCHAR NOT NULL,
    "action" VARCHAR NOT NULL,
    UNIQUE(user_type, resource, action)
);

-- users 
INSERT INTO permissions(user_type, resource, action) VALUES ('superadmin', 'users', 'create');
INSERT INTO permissions(user_type, resource, action) VALUES ('superadmin', 'users', 'update');
INSERT INTO permissions(user_type, resource, action) VALUES ('superadmin', 'users', 'delete');
INSERT INTO permissions(user_type, resource, action) VALUES ('superadmin', 'users', 'get');

INSERT INTO permissions(user_type, resource, action) VALUES ('user', 'users', 'update');
INSERT INTO permissions(user_type, resource, action) VALUES ('user', 'users', 'delete');
INSERT INTO permissions(user_type, resource, action) VALUES ('user', 'users', 'get');

-- categories
INSERT INTO permissions(user_type, resource, action) VALUES ('superadmin', 'categories', 'create');
INSERT INTO permissions(user_type, resource, action) VALUES ('superadmin', 'categories', 'delete');
INSERT INTO permissions(user_type, resource, action) VALUES ('superadmin', 'categories', 'update');

-- posts
INSERT INTO permissions(user_type, resource, action) VALUES ('superadmin', 'posts', 'create');
INSERT INTO permissions(user_type, resource, action) VALUES ('superadmin', 'posts', 'update');
INSERT INTO permissions(user_type, resource, action) VALUES ('superadmin', 'posts', 'delete');

INSERT INTO permissions(user_type, resource, action) VALUES ('user', 'posts', 'create');
INSERT INTO permissions(user_type, resource, action) VALUES ('user', 'posts', 'update');
INSERT INTO permissions(user_type, resource, action) VALUES ('user', 'posts', 'delete');

-- comments
INSERT INTO permissions(user_type, resource, action) VALUES ('superadmin', 'comments', 'create');
INSERT INTO permissions(user_type, resource, action) VALUES ('superadmin', 'comments', 'update');
INSERT INTO permissions(user_type, resource, action) VALUES ('superadmin', 'comments', 'delete');

INSERT INTO permissions(user_type, resource, action) VALUES ('user', 'comments', 'create');
INSERT INTO permissions(user_type, resource, action) VALUES ('user', 'comments', 'update');
INSERT INTO permissions(user_type, resource, action) VALUES ('user', 'comments', 'delete');

-- likes
INSERT INTO permissions(user_type, resource, action) VALUES ('user', 'likes', 'create');
INSERT INTO permissions(user_type, resource, action) VALUES ('user', 'likes', 'get');

INSERT INTO permissions(user_type, resource, action) VALUES ('superadmin', 'likes', 'create');
INSERT INTO permissions(user_type, resource, action) VALUES ('superadmin', 'likes', 'get');