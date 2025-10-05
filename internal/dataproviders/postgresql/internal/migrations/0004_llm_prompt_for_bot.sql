-- +goose Up

ALTER TABLE bots
ADD COLUMN llm_prompt TEXT;

-- +goose Down

ALTER TABLE bots
DROP COLUMN llm_prompt;