CREATE TABLE passwords (
                           id INTEGER PRIMARY KEY AUTOINCREMENT,
                           name VARCHAR(255) NOT NULL,
                           resource VARCHAR(255) NOT NULL,
                           password VARCHAR(255) NOT NULL,
                           salt TEXT,
                           iv TEXT
);
