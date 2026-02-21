-- Initialize all databases for microservices

-- Create user_db database
CREATE DATABASE user_db;
GRANT ALL PRIVILEGES ON DATABASE user_db TO postgres;

-- Create catalog_db database
CREATE DATABASE catalog_db;
GRANT ALL PRIVILEGES ON DATABASE catalog_db TO postgres;

-- Create order_db database
CREATE DATABASE order_db;
GRANT ALL PRIVILEGES ON DATABASE order_db TO postgres;

-- Create payment_db database
CREATE DATABASE payment_db;
GRANT ALL PRIVILEGES ON DATABASE payment_db TO postgres;

-- Create shipping_db database
CREATE DATABASE shipping_db;
GRANT ALL PRIVILEGES ON DATABASE shipping_db TO postgres;

-- Create notification_db database
CREATE DATABASE notification_db;
GRANT ALL PRIVILEGES ON DATABASE notification_db TO postgres;
