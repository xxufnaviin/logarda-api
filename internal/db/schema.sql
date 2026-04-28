CREATE TABLE users(
    username VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    awskeyID VARCHAR,
    awskeySecret VARCHAR,
    awsRegion VARCHAR,
    PRIMARY KEY (username)
);