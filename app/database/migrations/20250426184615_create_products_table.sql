-- +++ UP Migration
CREATE TABLE products (
 id BIGINT AUTO_INCREMENT PRIMARY KEY,
    reference VARCHAR(255) UNIQUE ,
    store_id BIGINT,
    category_id BIGINT,
    name VARCHAR(255),
    description TEXT,
    price DECIMAL(10, 2),
    margin DECIMAL(10, 2),
    stock INT,
    sold INT,
    images JSON,
    received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);
-- --- DOWN Migration
DROP TABLE IF EXISTS products;