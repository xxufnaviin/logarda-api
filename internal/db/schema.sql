CREATE TABLE metrics (
    metricTime TIMESTAMP NOT NULL,
    instanceID VARCHAR NOT NULL,
    cpu DOUBLE PRECISION NOT NULL,
    network DOUBLE PRECISION NOT NULL,
    memory DOUBLE PRECISION NOT NULL,
    username VARCHAR NOT NULL, 
    
    PRIMARY KEY (metricTime, instanceID) --composite primary key to prevent double entries at same timestamp for single instance
    CONSTRAINT fk_metrics_user
        FOREIGN KEY (username)
        REFERENCES users(username)
);

CREATE TABLE predicted_metrics (
    metricTime TIMESTAMP NOT NULL,
    instanceID VARCHAR NOT NULL,
    cpu DOUBLE PRECISION NOT NULL,
    network DOUBLE PRECISION NOT NULL,
    memory DOUBLE PRECISION NOT NULL,
    username VARCHAR NOT NULL, 
    
    PRIMARY KEY (metricTime, instanceID) --composite primary key to prevent double entries at same timestamp for single instance
    CONSTRAINT fk_predicted_metrics_user
        FOREIGN KEY (username)
        REFERENCES users(username)
);

CREATE TABLE logs (
    eventTime TIMESTAMP NOT NULL,
    errorCode VARCHAR NOT NULL,-- only store error logs
    errorMessage TEXT NOT NULL,
    serviceName VARCHAR NOT NULL,
    eventName VARCHAR NOT NULL,
    username VARCHAR NOT NULL,
    explanation VARCHAR,
    errorExplained BOOLEAN DEFAULT false,
    
    PRIMARY KEY (eventTime, errorCode, errorMessage) --composite primary keys to prevent double entries at same timestamp for same error code + error message
    CONSTRAINT fk_logs_user
        FOREIGN KEY (username)
        REFERENCES users(username)
);

CREATE TABLE users (
    username VARCHAR NOT NULL PRIMARY KEY,
    password VARCHAR NOT NULL,
    awskeyid VARCHAR,
    awskeysecret VARCHAR,
    awsregion VARCHAR,
    collector_on BOOLEAN DEFAULT false
);