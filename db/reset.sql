-- Create organization table
DROP TABLE IF EXISTS projects CASCADE;
DROP TABLE IF EXISTS organizations CASCADE;

-- Create organization table
CREATE TABLE organizations
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE
);

-- Create project table
CREATE TABLE projects
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE
);

CREATE TABLE accounts
(
    id              SERIAL PRIMARY KEY,
    user_uuid       UUID UNIQUE,
    email           VARCHAR(255) UNIQUE,
    password        VARCHAR(255),
    organization_id INTEGER REFERENCES organizations (id)
);


-- Insert into organization
INSERT INTO organizations (name)
VALUES ('Packlify');

-- Insert into project, associating with the organization and toolkit
INSERT INTO projects (name)
VALUES ('test');