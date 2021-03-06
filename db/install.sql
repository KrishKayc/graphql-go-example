-- CREATE DATABASE FOR OUR APP 'bookCab' --
CREATE DATABASE BookCab;

-- CREATE REGION TABLE --

USE BookCab;
CREATE TABLE Region (
    id BIGINT UNSIGNED AUTO_INCREMENT,
    name TEXT NOT NULL,
    minLatitude DECIMAL NOT NULL,
    maxLatitude DECIMAL NOT NULL,
    minLongitude DECIMAL NOT NULL,
    maxLongitude DECIMAL NOT NULL,
    PRIMARY KEY (id)

);
-- CREATE LOCATION TABLE --
CREATE TABLE Location (
    id BIGINT UNSIGNED AUTO_INCREMENT,
    name TEXT NOT NULL,
    latitude DECIMAL NOT NULL,
    longitude DECIMAL NOT NULL,
    PRIMARY KEY (id)
);

-- CREATE USER TABLE
CREATE TABLE User (
    id BIGINT UNSIGNED AUTO_INCREMENT,
    name TEXT NOT NULL,
    locationId BIGINT UNSIGNED,
    phoneNumber BIGINT,
    PRIMARY KEY (id),
    FOREIGN KEY (locationId) REFERENCES Location(id)
);

-- CREATE CAB TABLE
CREATE TABLE Cab (
    id BIGINT UNSIGNED AUTO_INCREMENT,
    type SMALLINT,
    driverId BIGINT UNSIGNED,
    number TEXT,
    PRIMARY KEY (id),
    FOREIGN KEY (driverId) REFERENCES User(id)
);

-- CREATE BOOKING TABLE

CREATE TABLE Booking(
    id BIGINT UNSIGNED AUTO_INCREMENT,
    cabId BIGINT UNSIGNED,
    userId BIGINT UNSIGNED,
    pickUpLocationId BIGINT UNSIGNED,
    dropLocationId BIGINT UNSIGNED,
    booked DATETIME,
    PRIMARY KEY (id),
    FOREIGN KEY (cabId) REFERENCES Cab(id),
    FOREIGN KEY (userId) REFERENCES User(id),
    FOREIGN KEY (pickUpLocationId) REFERENCES Location(id),
    FOREIGN KEY (dropLocationId) REFERENCES Location(id)
);