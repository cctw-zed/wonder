# Database Operations Guide

This document provides essential commands for managing the Wonder project database operations.

## Environment Information

- **Database**: PostgreSQL 15.14 (Homebrew)
- **Host**: localhost:5432
- **PostgreSQL Path**: `/opt/homebrew/opt/postgresql@15/bin/`

## Database Configuration

### Development Environment
- **Database**: `wonder_dev`
- **User**: `dev`
- **Password**: `dev`
- **Port**: 5432

### Testing Environment
- **Database**: `wonder_test`
- **User**: `test`
- **Password**: `test`
- **Port**: 5432

## Basic Connection Commands

### Connect to PostgreSQL as Superuser
```bash
# Connect to default postgres database
/opt/homebrew/opt/postgresql@15/bin/psql -h localhost -p 5432 -U $(whoami) -d postgres

# Alternative with environment PATH
psql -h localhost -p 5432 -U $(whoami) -d postgres
```

### Connect to Project Databases
```bash
# Development database
/opt/homebrew/opt/postgresql@15/bin/psql -h localhost -p 5432 -U dev -d wonder_dev

# Testing database
/opt/homebrew/opt/postgresql@15/bin/psql -h localhost -p 5432 -U test -d wonder_test
```

## Service Management

### PostgreSQL Service Control
```bash
# Start PostgreSQL service
brew services start postgresql@15

# Stop PostgreSQL service
brew services stop postgresql@15

# Restart PostgreSQL service
brew services restart postgresql@15

# Check service status
brew services list | grep postgres
```

### Check PostgreSQL Process
```bash
# Check if PostgreSQL is running
ps aux | grep postgres | grep -v grep

# Check PostgreSQL logs
tail -f /opt/homebrew/var/log/postgresql@15.log
```

## Database Information Commands

### List Databases
```sql
-- Connect as superuser first
\l
-- or
SELECT datname FROM pg_database;
```

### List Tables
```sql
-- After connecting to a specific database
\dt
-- or
SELECT tablename FROM pg_tables WHERE schemaname = 'public';
```

### Describe Table Structure
```sql
-- Show table structure
\d table_name
\d users

-- Show detailed table information
\d+ table_name
```

### List Users/Roles
```sql
-- List all users
\du
-- or
SELECT rolname FROM pg_roles;
```

## User and Database Management

### Create User
```sql
-- Create user with password and database creation privileges
CREATE USER username WITH PASSWORD 'password' CREATEDB;

-- Example
CREATE USER dev WITH PASSWORD 'dev' CREATEDB;
CREATE USER test WITH PASSWORD 'test' CREATEDB;
```

### Create Database
```sql
-- Create database with specific owner
CREATE DATABASE database_name OWNER username;

-- Examples
CREATE DATABASE wonder_dev OWNER dev;
CREATE DATABASE wonder_test OWNER test;
```

### Grant Permissions
```sql
-- Grant all privileges on database
GRANT ALL PRIVILEGES ON DATABASE database_name TO username;

-- Grant table permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO username;
```

### Drop Database/User (Caution!)
```sql
-- Drop database (WARNING: This deletes all data!)
DROP DATABASE database_name;

-- Drop user
DROP USER username;
```

## Table Operations

### Users Table Schema
```sql
-- Create users table (matches GORM structure)
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(64) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
```

### Auto-Update Trigger for updated_at
```sql
-- Create function for auto-updating updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

## Data Operations

### Insert Sample Data
```sql
-- Insert a test user
INSERT INTO users (id, email, name)
VALUES ('test-id-001', 'test@example.com', 'Test User');
```

### Query Data
```sql
-- Select all users
SELECT * FROM users;

-- Find user by email
SELECT * FROM users WHERE email = 'test@example.com';

-- Count total users
SELECT COUNT(*) FROM users;
```

### Update Data
```sql
-- Update user name (will automatically update updated_at via trigger)
UPDATE users SET name = 'Updated Name' WHERE email = 'test@example.com';
```

### Delete Data
```sql
-- Delete specific user
DELETE FROM users WHERE email = 'test@example.com';

-- Delete all users (CAUTION!)
DELETE FROM users;
```

## Backup and Restore

### Backup Database
```bash
# Backup single database
/opt/homebrew/opt/postgresql@15/bin/pg_dump -h localhost -p 5432 -U dev -d wonder_dev > wonder_dev_backup.sql

# Backup with custom format (recommended)
/opt/homebrew/opt/postgresql@15/bin/pg_dump -h localhost -p 5432 -U dev -d wonder_dev -Fc > wonder_dev_backup.dump
```

### Restore Database
```bash
# Restore from SQL file
/opt/homebrew/opt/postgresql@15/bin/psql -h localhost -p 5432 -U dev -d wonder_dev < wonder_dev_backup.sql

# Restore from custom format
/opt/homebrew/opt/postgresql@15/bin/pg_restore -h localhost -p 5432 -U dev -d wonder_dev wonder_dev_backup.dump
```

## Monitoring and Maintenance

### Check Database Size
```sql
-- Database sizes
SELECT datname, pg_size_pretty(pg_database_size(datname)) as size
FROM pg_database;

-- Table sizes
SELECT tablename, pg_size_pretty(pg_total_relation_size(tablename::regclass)) as size
FROM pg_tables WHERE schemaname = 'public';
```

### Active Connections
```sql
-- Show active connections
SELECT datname, usename, application_name, client_addr, state
FROM pg_stat_activity
WHERE state = 'active';
```

### Performance Monitoring
```sql
-- Show slow queries (if enabled)
SELECT query, mean_time, calls
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;
```

## Troubleshooting

### Common Issues

#### Connection Issues
```bash
# Check if PostgreSQL is running
brew services list | grep postgres

# Check port availability
lsof -i :5432

# Check PostgreSQL logs
tail -20 /opt/homebrew/var/log/postgresql@15.log
```

#### Permission Issues
```sql
-- Check user permissions
\du

-- Grant missing permissions
GRANT ALL PRIVILEGES ON DATABASE wonder_dev TO dev;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO dev;
```

#### Lock File Issues
```bash
# Remove lock files if PostgreSQL fails to start
sudo rm -f /tmp/.s.PGSQL.5432*
brew services restart postgresql@15
```

### Reset Everything (Nuclear Option)
```bash
# Stop PostgreSQL
brew services stop postgresql@15

# Remove data directory (WARNING: Deletes all data!)
rm -rf /opt/homebrew/var/postgresql@15

# Reinitialize database
/opt/homebrew/opt/postgresql@15/bin/initdb -D /opt/homebrew/var/postgresql@15

# Start service
brew services start postgresql@15

# Recreate users and databases using commands above
```

## Environment Variables

### For Application Use
```bash
# Development
export DB_HOST=localhost
export DB_PORT=5432
export DB_USERNAME=dev
export DB_PASSWORD=dev
export DB_DATABASE=wonder_dev

# Testing
export DB_HOST=localhost
export DB_PORT=5432
export DB_USERNAME=test
export DB_PASSWORD=test
export DB_DATABASE=wonder_test
```

### For Scripts
```bash
# Set PostgreSQL binary path
export PATH="/opt/homebrew/opt/postgresql@15/bin:$PATH"

# Connection string format
export DATABASE_URL="postgres://username:password@localhost:5432/database_name"
```

## Quick Reference

### Essential Commands Summary
```bash
# Service
brew services start postgresql@15
brew services stop postgresql@15

# Connect
psql -h localhost -p 5432 -U dev -d wonder_dev

# Backup
pg_dump -h localhost -p 5432 -U dev -d wonder_dev > backup.sql

# Restore
psql -h localhost -p 5432 -U dev -d wonder_dev < backup.sql
```

### Essential SQL Commands
```sql
-- Information
\l          -- List databases
\dt         -- List tables
\d table    -- Describe table
\du         -- List users

-- Data
SELECT * FROM users;
INSERT INTO users (id, email, name) VALUES ('id', 'email', 'name');
UPDATE users SET name = 'new_name' WHERE id = 'user_id';
DELETE FROM users WHERE id = 'user_id';
```

---

**Note**: Always backup your data before performing destructive operations. This guide assumes PostgreSQL 15 installed via Homebrew on macOS.

**Last Updated**: 2025-09-23
**Project**: Wonder
**Environment**: Development/Testing