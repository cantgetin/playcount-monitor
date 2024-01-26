-- +migrate Up
CREATE TABLE tracks
(
    id                         serial primary key,
    tracked_at                 timestamp default NOW()
);

INSERT INTO tracks (id, tracked_at)
VALUES (1, '2022-10-10 11:30:30');

-- +migrate Down

DROP TABLE tracks