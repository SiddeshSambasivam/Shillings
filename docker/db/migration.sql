CREATE TABLE users (
    user_id int(11) NOT NULL AUTO_INCREMENT,
    first_name varchar(255) NOT NULL,
    middle_name varchar(255) default "",
    last_name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    phone varchar(255) default "",
    balance FLOAT(6) NOT NULL DEFAULT 0,
    created_at int NOT NULL,
    updated_at int NOT NULL,
    PRIMARY KEY (user_id)    
) ENGINE=InnoDB;

ALTER TABLE users AUTO_INCREMENT=10001;
CREATE INDEX idx_email on users(email);

-- INSERT INTO users(first_name, last_name, email, balance, created_at, updated_at) VALUES
-- ('John', 'Doe', 'JohnDoe@test.com', 0, 1518098983, 1518098983),
-- ('Jack', 'Sparrow', 'JackSparrow@pirates.com', 0, 1518098983, 1518098983);

CREATE TABLE transactions (
    transaction_id int(11) NOT NULL AUTO_INCREMENT,
    sender_id int(11) NOT NULL,
    receiver_id int(11) NOT NULL,
    amount FLOAT(6) NOT NULL,
    created_at int NOT NULL,        
    PRIMARY KEY (transaction_id)    
) ENGINE= InnoDB;

CREATE INDEX idx_email on transactions(sender_id, receiver_id);

ALTER TABLE transactions AUTO_INCREMENT=10001;

-- INSERT INTO transactions(sender_id, receiver_id, amount, created_at) VALUES
-- (10002, 10003, 10, 1518098983),
-- (10001, 10002, 20, 1518098983);

CREATE TABLE credentials (
    credential_id int(11) NOT NULL AUTO_INCREMENT,
    user_id int(11) NOT NULL,    
    email varchar(255) NOT NULL,
    password varchar(255) NOT NULL,
    created_at int NOT NULL,
    updated_at int NOT NULL,
    PRIMARY KEY (credential_id)    
) ENGINE= InnoDB;

ALTER TABLE credentials AUTO_INCREMENT=10001;
CREATE INDEX idx_email on credentials(email);

-- INSERT INTO credentials(user_id, salt, email, password, created_at) VALUES
-- (10002, 'salt', 'JohnDoe@test.com', 'password', 1518098983);