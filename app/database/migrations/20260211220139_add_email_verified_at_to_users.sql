-- +++ UP Migration
ALTER TABLE users ADD COLUMN email_verified_at TIMESTAMP NULL DEFAULT NULL AFTER email;
-- --- DOWN Migration
ALTER TABLE users DROP COLUMN email_verified_at;
