-- +goose Up

CREATE TABLE bots(
        id SERIAL8 PRIMARY KEY,
        name TEXT NOT NULL,
        bot_tag TEXT NOT NULL,
        token TEXT NOT NULL,
        enabled BOOLEAN NOT NULL,
        description TEXT,
        emoji_list TEXT ARRAY,
        emoji_chance FLOAT,
        tags TEXT ARRAY,
        allowed_chats INT8 ARRAY,
        create_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        update_at TIMESTAMPTZ
);

CREATE TABLE
    "quotes"(
        id SERIAL8 PRIMARY KEY,
        bot_id INT8 NOT NULL REFERENCES bots (id) ON UPDATE CASCADE ON DELETE CASCADE,
        "text" TEXT NOT NULL,
        user_id INT8,
        chat_id INT8,
        create_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        update_at TIMESTAMPTZ
    );

CREATE TABLE
    moderators(
        user_id INT8 NOT NULL,
        bot_id INT8 NOT NULL REFERENCES bots (id) ON UPDATE CASCADE ON DELETE CASCADE,
        create_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        description TEXT,
        PRIMARY KEY(user_id, bot_id)
    );