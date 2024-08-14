CREATE TABLE users
(
    id        VARCHAR PRIMARY KEY,
    user_name VARCHAR UNIQUE NOT NULL
);
CREATE TABLE user_data
(
    user_id VARCHAR PRIMARY KEY REFERENCES users (id),
    name    VARCHAR,
    surname VARCHAR,
    role    VARCHAR
);