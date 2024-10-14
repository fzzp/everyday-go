CREATE TABLE products (
    id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
    productName varchar(255) NOT NULL,
    productPrice int NOT NULL
);

CREATE TABLE category (
    id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
    categoryName varchar(255) NOT NULL UNIQUE
);

CREATE TABLE category_products(
    categoryId int NOT NULL,
    productId int NOT NULL
);