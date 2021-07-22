
-- +migrate Up
CREATE SEQUENCE events_id_seq;

CREATE TABLE IF NOT EXISTS events (
    eventid           INTEGER       NOT NULL DEFAULT nextval('events_id_seq'),
    event_name        VARCHAR(255)  NOT NULL,
    event_start       DATE          NOT NULL,
    event_end         DATE          NOT NULL,
    event_description VARCHAR(255),
    event_alert       DATE,
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
DROP SEQUENCE IF EXISTS events_id_seq;
