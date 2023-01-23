-- name: CreateArtist :one
INSERT INTO artists 
    (name, image) 
VALUES 
    ($1, $2) RETURNING id;

-- name: CreateAlbum :one
INSERT INTO albums 
    (artist_id, name, cover) 
VALUES
    ($1, $2, $3) RETURNING id;

-- name: CreateSong :exec
INSERT INTO songs 
    (album_id, artist_id, title, length, track, path, mtime) 
VALUES 
    ($1, $2, $3, $4, $5, $6, $7);