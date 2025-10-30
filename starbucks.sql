CREATE TABLE  menu_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_menu_name_not_empty CHECK (TRIM(name) != ''),
    CONSTRAINT chk_menu_price_positive CHECK (price > 0),
    CONSTRAINT uq_menu_name UNIQUE (name)
);

DROP TABLE menu_items ;

INSERT INTO menu_items (id, name, price) VALUES
(1, 'Caffè Americano (Tall)', 39000),
(2, 'Caffè Latte (Tall)', 45000),
(3, 'Cappuccino (Tall)', 45000),
(4, 'Caramel Macchiato (Tall)', 55000),
(5, 'Espresso (Double Shot)', 35000),
(6, 'Mocha Frappuccino', 58000),
(7, 'Java Chip Frappuccino', 60000),
(8, 'Green Tea Latte', 55000),
(9, 'Signature Chocolate', 52000),
(10, 'Vanilla Sweet Cream Cold Brew', 56000);

SELECT * FROM menu_items;

CREATE TABLE  orders (
    id SERIAL PRIMARY KEY,
    total INT NOT NULL,
    order_date VARCHAR(50) NOT NULL,
    order_time VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_order_total_positive CHECK (total > 0),
    CONSTRAINT chk_order_date_not_empty CHECK (TRIM(order_date) != ''),
    CONSTRAINT chk_order_time_not_empty CHECK (TRIM(order_time) != '')
);

CREATE TABLE  order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    menu_id VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT chk_order_item_menu_id_not_empty CHECK (TRIM(menu_id) != ''),
    CONSTRAINT chk_order_item_name_not_empty CHECK (TRIM(name) != ''),
    CONSTRAINT chk_order_item_price_positive CHECK (price > 0),
    CONSTRAINT chk_order_item_quantity_positive CHECK (quantity > 0)
);

SELECT* FROM orders;
SELECT* FROM order_items;
