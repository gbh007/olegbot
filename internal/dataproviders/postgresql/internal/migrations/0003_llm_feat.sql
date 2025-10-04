-- +goose Up

ALTER TABLE bots
ADD COLUMN llm_chance FLOAT,
ADD COLUMN llm_allowed_chats INT8 ARRAY;

-- +goose Down

ALTER TABLE bots
DROP COLUMN llm_chance,
DROP COLUMN llm_allowed_chats;