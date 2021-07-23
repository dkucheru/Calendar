
-- +migrate Up
CREATE TABLE IF NOT EXISTS events (
    eventid           BIGSERIAL                     NOT NULL,
    event_name        VARCHAR(255)                  NOT NULL,
    event_start       TIMESTAMP WITHOUT TIME ZONE   NOT NULL,
    event_end         TIMESTAMP WITHOUT TIME ZONE   NOT NULL,
    event_description VARCHAR(255),
    event_alert       TIMESTAMP WITHOUT TIME ZONE,
    PRIMARY KEY (eventid)
);

CREATE TABLE IF NOT EXISTS users (
    username     VARCHAR(255)                NOT NULL,
    hashedpass   VARCHAR(255)                NOT NULL,
    userlocation VARCHAR(255)                NOT NULL,
    PRIMARY KEY (username)
);

-- +migrate Down
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS events;
