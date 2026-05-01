-- Очистити таблиці
TRUNCATE TABLE order_items, orders, cart_items, products, categories, users CASCADE;

-- ============================================
-- Користувачі (пароль для всіх: Test1234!)
-- ============================================
INSERT INTO users (id, email, password_hash, first_name, last_name, role) VALUES
('11111111-1111-1111-1111-111111111111', 'admin@test.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
 'Admin', 'User', 'admin'),
('22222222-2222-2222-2222-222222222222', 'user@test.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
 'John', 'Doe', 'customer');

-- ============================================
-- Категорії
-- ============================================
INSERT INTO categories (id, name, description) VALUES
('c1111111-1111-1111-1111-111111111111', 'Electronics', 'Electronic devices'),
('c2222222-2222-2222-2222-222222222222', 'Books', 'Books and e-books'),
('c3333333-3333-3333-3333-333333333333', 'Clothing', 'Fashion and apparel');

-- ============================================
-- Товари (мінімум 20 на категорію)
-- ============================================

-- ELECTRONICS (20)
INSERT INTO products (name, description, price, stock, category_id) VALUES
('Laptop Pro', 'High-performance laptop', 1299.99, 50, 'c1111111-1111-1111-1111-111111111111'),
('Wireless Mouse', 'Ergonomic mouse', 29.99, 200, 'c1111111-1111-1111-1111-111111111111'),
('USB-C Hub', '7-in-1 hub', 49.99, 150, 'c1111111-1111-1111-1111-111111111111'),
('Keyboard', 'Mechanical keyboard', 89.99, 100, 'c1111111-1111-1111-1111-111111111111'),
('Webcam', '1080p webcam', 59.99, 80, 'c1111111-1111-1111-1111-111111111111'),
('Headphones', 'Wireless headphones', 149.99, 120, 'c1111111-1111-1111-1111-111111111111'),
('SSD 1TB', 'Portable SSD', 119.99, 90, 'c1111111-1111-1111-1111-111111111111'),
('Monitor 27"', '4K monitor', 399.99, 60, 'c1111111-1111-1111-1111-111111111111'),
('Bluetooth Speaker', 'Portable speaker', 79.99, 140, 'c1111111-1111-1111-1111-111111111111'),
('Smartwatch', 'Fitness tracking watch', 199.99, 110, 'c1111111-1111-1111-1111-111111111111'),
('Gaming Mouse', 'RGB gaming mouse', 49.99, 180, 'c1111111-1111-1111-1111-111111111111'),
('Gaming Keyboard', 'RGB mechanical keyboard', 129.99, 90, 'c1111111-1111-1111-1111-111111111111'),
('Noise Cancelling Headphones', 'Premium ANC headphones', 249.99, 70, 'c1111111-1111-1111-1111-111111111111'),
('Portable Charger', '20000mAh power bank', 39.99, 300, 'c1111111-1111-1111-1111-111111111111'),
('Smartphone Stand', 'Adjustable stand', 14.99, 500, 'c1111111-1111-1111-1111-111111111111'),
('Tablet 10"', 'Android tablet', 229.99, 85, 'c1111111-1111-1111-1111-111111111111'),
('Wireless Charger', 'Fast charging pad', 24.99, 250, 'c1111111-1111-1111-1111-111111111111'),
('Action Camera', '4K action camera', 159.99, 75, 'c1111111-1111-1111-1111-111111111111'),
('Drone Mini', 'Compact drone', 299.99, 40, 'c1111111-1111-1111-1111-111111111111'),
('VR Headset', 'Virtual reality headset', 349.99, 35, 'c1111111-1111-1111-1111-111111111111'),

-- BOOKS (20)
('Go Programming', 'Learn Go', 39.99, 300, 'c2222222-2222-2222-2222-222222222222'),
('Clean Architecture', 'Design patterns', 44.99, 250, 'c2222222-2222-2222-2222-222222222222'),
('Docker Deep Dive', 'Container guide', 49.99, 200, 'c2222222-2222-2222-2222-222222222222'),
('PostgreSQL Guide', 'Database book', 54.99, 180, 'c2222222-2222-2222-2222-222222222222'),
('System Design', 'Interview prep', 42.99, 220, 'c2222222-2222-2222-2222-222222222222'),
('Microservices Patterns', 'Service architecture', 59.99, 160, 'c2222222-2222-2222-2222-222222222222'),
('Kubernetes Up & Running', 'K8s guide', 49.99, 190, 'c2222222-2222-2222-2222-222222222222'),
('Designing Data-Intensive Apps', 'Distributed systems', 64.99, 170, 'c2222222-2222-2222-2222-222222222222'),
('Algorithms', 'Computer science', 39.99, 210, 'c2222222-2222-2222-2222-222222222222'),
('Operating Systems', 'OS concepts', 69.99, 140, 'c2222222-2222-2222-2222-222222222222'),
('Networking Basics', 'Network fundamentals', 34.99, 260, 'c2222222-2222-2222-2222-222222222222'),
('AI Fundamentals', 'Artificial intelligence', 59.99, 150, 'c2222222-2222-2222-2222-222222222222'),
('Machine Learning', 'ML concepts', 74.99, 130, 'c2222222-2222-2222-2222-222222222222'),
('Deep Learning', 'Neural networks', 79.99, 120, 'c2222222-2222-2222-2222-222222222222'),
('Cybersecurity Essentials', 'Security basics', 44.99, 200, 'c2222222-2222-2222-2222-222222222222'),
('Cloud Computing', 'Cloud fundamentals', 39.99, 230, 'c2222222-2222-2222-2222-222222222222'),
('API Design', 'REST & gRPC', 49.99, 180, 'c2222222-2222-2222-2222-222222222222'),
('DevOps Handbook', 'DevOps guide', 54.99, 160, 'c2222222-2222-2222-2222-222222222222'),
('Scrum Guide', 'Agile methodology', 24.99, 300, 'c2222222-2222-2222-2222-222222222222'),
('Linux Bible', 'Linux administration', 69.99, 140, 'c2222222-2222-2222-2222-222222222222'),

-- CLOTHING (20)
('T-Shirt', 'Cotton t-shirt', 19.99, 500, 'c3333333-3333-3333-3333-333333333333'),
('Jeans', 'Blue jeans', 59.99, 300, 'c3333333-3333-3333-3333-333333333333'),
('Hoodie', 'Comfortable hoodie', 39.99, 200, 'c3333333-3333-3333-3333-333333333333'),
('Sneakers', 'Sport shoes', 79.99, 150, 'c3333333-3333-3333-3333-333333333333'),
('Jacket', 'Winter jacket', 129.99, 100, 'c3333333-3333-3333-3333-333333333333'),
('Cap', 'Baseball cap', 14.99, 400, 'c3333333-3333-3333-3333-333333333333'),
('Socks Pack', '5 pairs', 9.99, 600, 'c3333333-3333-3333-3333-333333333333'),
('Shorts', 'Summer shorts', 24.99, 350, 'c3333333-3333-3333-3333-333333333333'),
('Sweatpants', 'Comfortable pants', 34.99, 250, 'c3333333-3333-3333-3333-333333333333'),
('Boots', 'Leather boots', 99.99, 120, 'c3333333-3333-3333-3333-333333333333'),
('Scarf', 'Winter scarf', 12.99, 500, 'c3333333-3333-3333-3333-333333333333'),
('Gloves', 'Warm gloves', 9.99, 450, 'c3333333-3333-3333-3333-333333333333'),
('Belt', 'Leather belt', 19.99, 300, 'c3333333-3333-3333-3333-333333333333'),
('Sunglasses', 'UV protection', 29.99, 200, 'c3333333-3333-3333-3333-333333333333'),
('Backpack', 'Casual backpack', 49.99, 180, 'c3333333-3333-3333-3333-333333333333'),
('Dress', 'Summer dress', 39.99, 160, 'c3333333-3333-3333-3333-333333333333'),
('Skirt', 'Pleated skirt', 29.99, 140, 'c3333333-3333-3333-3333-333333333333'),
('Formal Shirt', 'Office wear', 34.99, 220, 'c3333333-3333-3333-3333-333333333333'),
('Suit Jacket', 'Formal jacket', 89.99, 90, 'c3333333-3333-3333-3333-333333333333'),
('Running Shoes', 'Lightweight shoes', 69.99, 170, 'c3333333-3333-3333-3333-333333333333');

SELECT 'Seed data loaded!' AS status;
