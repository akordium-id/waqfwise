-- WaqfWise Community Edition - Rollback Initial Schema

DROP TABLE IF EXISTS payment_logs;
DROP TABLE IF EXISTS donations;
DROP TABLE IF EXISTS campaigns;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "uuid-ossp";
