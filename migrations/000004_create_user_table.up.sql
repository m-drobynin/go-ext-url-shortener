CREATE TABLE users (
    user_id UUID PRIMARY KEY
);

ALTER TABLE urls ADD COLUMN user_id UUID REFERENCES users (user_id);
