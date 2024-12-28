-- +goose Up

CREATE TABLE
    "stickers"(
        id SERIAL8 PRIMARY KEY,
        bot_id INT8 NOT NULL REFERENCES bots (id) ON UPDATE CASCADE ON DELETE CASCADE,
        file_id TEXT NOT NULL,
        user_id INT8,
        chat_id INT8,
        create_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        update_at TIMESTAMPTZ
    );

CREATE TABLE
    "gifs"(
        id SERIAL8 PRIMARY KEY,
        bot_id INT8 NOT NULL REFERENCES bots (id) ON UPDATE CASCADE ON DELETE CASCADE,
        file_id TEXT NOT NULL,
        user_id INT8,
        chat_id INT8,
        create_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        update_at TIMESTAMPTZ
    );


ALTER TABLE bots
ADD COLUMN sticker_chance FLOAT,
ADD COLUMN gif_chance FLOAT;

-- +goose Down

ALTER TABLE bots
DROP COLUMN sticker_chance,
DROP COLUMN gif_chance;