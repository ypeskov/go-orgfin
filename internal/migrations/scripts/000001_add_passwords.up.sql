CREATE TABLE users
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    email         VARCHAR(255) NOT NULL UNIQUE,
    name          VARCHAR(255) NOT NULL,
    hash_password VARCHAR(255) NOT NULL,
    last_login    DATETIME,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE encrypted_passwords
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER      NOT NULL,
    name       VARCHAR(255) NOT NULL,
    resource   VARCHAR(255) NOT NULL,
    login      VARCHAR(255),
    password   VARCHAR(255) NOT NULL,
    salt       TEXT,
    iv         TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);


CREATE TABLE encrypted_credit_cards
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id        INTEGER      NOT NULL,
    card_number    TEXT         NOT NULL,
    cardholder     VARCHAR(255) NOT NULL,
    expiration     VARCHAR(7)   NOT NULL, -- Format MM/YYYY
    cvv            TEXT         NOT NULL,
    encrypted_data TEXT         NOT NULL,
    salt           TEXT,
    iv             TEXT,
    created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE encrypted_notes
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id        INTEGER      NOT NULL,
    title          VARCHAR(255) NOT NULL,
    encrypted_note TEXT         NOT NULL,
    salt           TEXT,
    iv             TEXT,
    created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
