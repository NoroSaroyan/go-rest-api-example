-- Create the todos table
CREATE TABLE IF NOT EXISTS todos
(
    id         SERIAL PRIMARY KEY,
    title      TEXT      NOT NULL,
    completed  BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index for filtering by completed status
CREATE INDEX IF NOT EXISTS idx_todos_completed ON todos (completed);
