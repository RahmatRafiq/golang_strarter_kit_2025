-- +++ UP Migration
ALTER TABLE stores 
ADD COLUMN new_column_name VARCHAR(255) NOT NULL DEFAULT 'default_value',
ADD COLUMN another_column_name INT NOT NULL DEFAULT 0,
ADD COLUMN yet_another_column_name TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- --- DOWN Migration
ALTER TABLE stores 
DROP COLUMN new_column_name,
DROP COLUMN another_column_name,
DROP COLUMN yet_another_column_name;
