CREATE DATABASE IF NOT EXISTS Library;
USE Library;

CREATE TABLE IF NOT EXISTS User (
    UserID INT AUTO_INCREMENT PRIMARY KEY,
    Username VARCHAR(255) NOT NULL UNIQUE,
    Pass VARCHAR(255) NOT NULL,
    Role VARCHAR(255) NOT NULL,
    Created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    AdminRequest ENUM('NoRequest', 'Pending', 'Accepted', 'Denied') DEFAULT 'NoRequest'
);

CREATE TABLE IF NOT EXISTS Books (
    BookID INT AUTO_INCREMENT PRIMARY KEY,
    Title VARCHAR(255) NOT NULL,
    Author VARCHAR(255),
    Genre VARCHAR(255),
    Quantity INT 
);


CREATE TABLE IF NOT EXISTS BookRequests (
    RequestID INT AUTO_INCREMENT PRIMARY KEY,
    UserID INT,
    BookID INT,
    RequestDate TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    AcceptDate TIMESTAMP DEFAULT NULL,
    Status ENUM('Pending', 'Accepted', 'Returned', 'Denied') DEFAULT 'Pending',
    FOREIGN KEY (UserID) REFERENCES User(UserID),
    FOREIGN KEY (BookID) REFERENCES Books(BookID)
);

