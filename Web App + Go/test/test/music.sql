DROP TABLE IF EXISTS musician;

CREATE TABLE musician (
    music_id serial PRIMARY KEY,
    full_name text NOT NULL,
    album text NOT NULL,
    genre text NOT NULL,
    date_released date NOT NULL,
    artist text NOT NULL
)