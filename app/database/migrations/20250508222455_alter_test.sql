-- +++ UP Migration
ALTER TABLE test 
ADD COLUMN new_column_name VARCHAR(255) DEFAULT 'default_value';

-- --- DOWN Migration
ALTER TABLE test 
DROP COLUMN new_column_name;
