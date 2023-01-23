// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: music.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createAlbum = `-- name: CreateAlbum :one
INSERT INTO albums 
    (artist_id, name, cover) 
VALUES
    ($1, $2, $3) RETURNING id
`

type CreateAlbumParams struct {
	ArtistID int64
	Name     string
	Cover    string
}

func (q *Queries) CreateAlbum(ctx context.Context, arg CreateAlbumParams) (int64, error) {
	row := q.db.QueryRow(ctx, createAlbum, arg.ArtistID, arg.Name, arg.Cover)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createArtist = `-- name: CreateArtist :one
INSERT INTO artists 
    (name, image) 
VALUES 
    ($1, $2) RETURNING id
`

type CreateArtistParams struct {
	Name  string
	Image sql.NullString
}

func (q *Queries) CreateArtist(ctx context.Context, arg CreateArtistParams) (int64, error) {
	row := q.db.QueryRow(ctx, createArtist, arg.Name, arg.Image)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createSong = `-- name: CreateSong :exec
INSERT INTO songs 
    (album_id, artist_id, title, length, track, path, mtime) 
VALUES 
    ($1, $2, $3, $4, $5, $6, $7)
`

type CreateSongParams struct {
	AlbumID  int64
	ArtistID int64
	Title    string
	Length   float64
	Track    sql.NullInt32
	Path     string
	Mtime    int32
}

func (q *Queries) CreateSong(ctx context.Context, arg CreateSongParams) error {
	_, err := q.db.Exec(ctx, createSong,
		arg.AlbumID,
		arg.ArtistID,
		arg.Title,
		arg.Length,
		arg.Track,
		arg.Path,
		arg.Mtime,
	)
	return err
}
