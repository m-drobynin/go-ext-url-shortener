CREATE TABLE users (
    user_id UUID PRIMARY KEY
);

ALTER TABLE urls ADD FOREIGN KEY (user_id) REFERENCES users (user_id);
