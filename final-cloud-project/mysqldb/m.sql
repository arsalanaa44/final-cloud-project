CREATE DATABASE bepa;

USE bepa;

CREATE TABLE prices (
        coin_name VARCHAR(255) NOT NULL,
        timestamp DATETIME NOT NULL,
        price DECIMAL(10, 2) NOT NULL
    );

CREATE TABLE alertSubscription (
        user_email VARCHAR(255) NOT NULL,
        coin_name VARCHAR(255) NOT NULL,
        difference_percentage DECIMAL(5, 2) NOT NULL
    );
    
select * from alertSubscription ;
select * from prices ;
