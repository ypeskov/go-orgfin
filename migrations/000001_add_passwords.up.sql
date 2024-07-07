CREATE TABLE passwords (
                           id SERIAL PRIMARY KEY,
                           name VARCHAR(255) NOT NULL,
                           url VARCHAR(255) NOT NULL,
                           password VARCHAR(255) NOT NULL
);
