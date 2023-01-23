-- migrate:up
CREATE TABLE "artists"(
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "image" VARCHAR(255) NULL
);
ALTER TABLE
    "artists" ADD CONSTRAINT "artists_name_unique" UNIQUE("name");

CREATE TABLE "albums"(
    "id" BIGSERIAL PRIMARY KEY,
    "artist_id" BIGSERIAL NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "cover" VARCHAR(255) NOT NULL
);

CREATE TABLE "songs"(
    "id" BIGSERIAL PRIMARY KEY,
    "album_id" BIGSERIAL NOT NULL,
    "artist_id" BIGSERIAL NOT NULL,
    "title" VARCHAR(255) NOT NULL,
    "length" DOUBLE PRECISION NOT NULL,
    "track" INTEGER NULL,
    "path" TEXT NOT NULL,
    "mtime" INTEGER NOT NULL
);
ALTER TABLE
    "albums" ADD CONSTRAINT "albums_artist_id_foreign" FOREIGN KEY("artist_id") REFERENCES "artists"("id");
ALTER TABLE
    "songs" ADD CONSTRAINT "songs_album_id_foreign" FOREIGN KEY("album_id") REFERENCES "albums"("id");
ALTER TABLE
    "songs" ADD CONSTRAINT "songs_artist_id_foreign" FOREIGN KEY("artist_id") REFERENCES "artists"("id");

-- migrate:down
DROP TABLE "songs";
DROP TABLE "albums";
DROP TABLE "artists";
