CREATE TYPE status_enum AS ENUM ('pending','cancelled','completed');


CREATE TABLE IF NOT EXISTS users(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) UNIQUE NOT NULL,
    password varchar,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at bigint DEFAULT 0
);

CREATE TABLE IF NOT EXISTS restaurants(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    address VARCHAR(100) NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    description text,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at bigint DEFAULT 0
);

CREATE TABLE IF NOT EXISTS reservations(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    restaurant_id uuid REFERENCES restaurants(id),
    reservation_time TIMESTAMP NOT NULL,
    status status_enum,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at bigint DEFAULT 0
);

CREATE TABLE IF NOT EXISTS menu(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurants_id uuid REFERENCES restaurants(id),
    name VARCHAR(100) NOT NULL,
    description text,
    price DECIMAL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at bigint DEFAULT 0
);

CREATE TABLE IF NOT EXISTS reservationsorders(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),  
    reservation_id uuid REFERENCES reservations(id),
    menu_item_id uuid REFERENCES menu(id),
    quantity INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at bigint DEFAULT 0
);