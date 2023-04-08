-- Instantiate an empty database.
CREATE DATABASE "rest-api-framework";
\connect "rest-api-framework";
CREATE SCHEMA IF NOT EXISTS "public";

-- Create a new full-privilege user.
CREATE USER "rest-api-framework-user" WITH PASSWORD 'secret123';
GRANT ALL ON DATABASE "rest-api-framework" TO "rest-api-framework-user";
GRANT USAGE ON SCHEMA "public" TO "rest-api-framework-user";
GRANT ALL ON ALL TABLES IN SCHEMA "public" TO "rest-api-framework-user";
ALTER DEFAULT PRIVILEGES IN SCHEMA "public" GRANT ALL ON TABLES TO "rest-api-framework-user";

-- Create a new readonly user.
CREATE USER "rest-api-framework-readonly-user" WITH PASSWORD 'secret123';
GRANT CONNECT ON DATABASE "rest-api-framework" TO "rest-api-framework-readonly-user";
GRANT USAGE ON SCHEMA "public" TO "rest-api-framework-readonly-user";
GRANT SELECT ON ALL TABLES IN SCHEMA "public" TO "rest-api-framework-readonly-user";
ALTER DEFAULT PRIVILEGES IN SCHEMA "public" GRANT SELECT ON TABLES TO "rest-api-framework-readonly-user";

-- Create version table that would manage migrations.
CREATE TABLE versions (
    id VARCHAR(255) PRIMARY KEY
);
