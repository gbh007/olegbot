CREATE TABLE
    "quotes"(
        id SERIAL8 PRIMARY KEY,
        "text" TEXT NOT NULL,
        create_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        user_id INT8,
        chat_id INT8
    );

CREATE TABLE
    moderators(
        user_id INT8 NOT NULL PRIMARY KEY,
        create_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        description TEXT
    );