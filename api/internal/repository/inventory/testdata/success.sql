INSERT INTO users(id, name, email, password, status)
VALUES
    (14753001,'Test User','test@example.com', 'password123', 'ACTIVE'),
    (14753002,'Test User2','test2@example.com', 'password@123', 'ACTIVE');

INSERT INTO products(id, name, description, status, price, stock)
VALUES
    (14753010, 'Test Product', 'test', 'ACTIVE', 2000, 100);

INSERT INTO orders(id, user_id, status, total_cost)
VALUES
    (14753010, 14753001, 'PENDING', 20),
    (14753011, 14753001, 'PENDING', 10);
