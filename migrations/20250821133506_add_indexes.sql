-- +goose Up
-- Index for fast user lookups
CREATE INDEX idx_users_email ON users(email);

-- Index for expense filtering
CREATE INDEX idx_expenses_status ON expenses(status);
CREATE INDEX idx_expenses_category ON expenses(category);

-- Index for expense report ownership
CREATE INDEX idx_reports_user_id ON expense_reports(user_id);

-- +goose Down
DROP INDEX idx_users_email;
DROP INDEX idx_expenses_status;
DROP INDEX idx_expenses_category;
DROP INDEX idx_reports_user_id;
