-- Drop triggers first (depend on tables and function)
DROP TRIGGER IF EXISTS update_collections_updated_at ON collections;
DROP TRIGGER IF EXISTS update_tags_updated_at ON tags;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON tokens;
DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_resources_updated_at ON resources;

-- Drop tables (dependents first, then parent tables)
DROP TABLE IF EXISTS resource_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS resources;
DROP TABLE IF EXISTS collections;
DROP TABLE IF EXISTS users;

-- Drop function last (nothing depends on it anymore)
DROP FUNCTION IF EXISTS update_updated_at_column();
