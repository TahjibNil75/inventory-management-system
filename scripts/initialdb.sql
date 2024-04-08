
-- Check if the database exists, and create it if it doesn't
SELECT 'CREATE DATABASE inventory'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'inventory');


