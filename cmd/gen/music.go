package gen

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/joeychilson/testdb/db"
	"github.com/joeychilson/testdb/db/sqlc"
	"github.com/joeychilson/testdb/gen"
	"github.com/spf13/cobra"
)

type MusicGenCmd struct {
	db *db.Postgres
}

func NewMusicGen(db *db.Postgres) *MusicGenCmd {
	return &MusicGenCmd{db: db}
}

func (g *MusicGenCmd) MusicGenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "music [flags]",
		Short: "Generate realistic fake data for the music schema",
		RunE: func(cmd *cobra.Command, args []string) error {
			artists, err := cmd.Flags().GetInt("artists")
			if err != nil {
				return err
			}

			albums, err := cmd.Flags().GetInt("albums")
			if err != nil {
				return err
			}

			songs, err := cmd.Flags().GetInt("songs")
			if err != nil {
				return err
			}

			return g.handleCmd(cmd.Context(), artists, albums, songs)
		},
	}

	cmd.Flags().IntP("artists", "r", 1, "The number of artists to generate")
	cmd.Flags().IntP("albums", "a", 1, "The number of albums per artists to generate")
	cmd.Flags().IntP("songs", "s", 1, "The number of songs per albums to generate")
	return cmd
}

func (g *MusicGenCmd) handleCmd(ctx context.Context, artists int, albums int, songs int) error {
	log.Printf("Generating %d artists, %d albums, and %d songs for the music schema", artists, albums, songs)

	for i := 0; i < artists; i++ {
		gofakeit.Seed(0)

		var genArtist *gen.Artist
		gofakeit.Struct(&genArtist)

		artist := sqlc.CreateArtistParams{
			Name:  genArtist.Name,
			Image: sql.NullString{String: genArtist.Image, Valid: genArtist.Image != ""},
		}

		artistsID, err := g.db.CreateArtist(ctx, artist)
		if err != nil {
			log.Fatalf("failed to create artist: %v", err)
		}

		log.Printf("Created artist %s", genArtist.Name)

		for i := 0; i < albums; i++ {
			var genAlbum *gen.Album
			gofakeit.Struct(&genAlbum)

			album := sqlc.CreateAlbumParams{
				ArtistID: artistsID,
				Name:     strings.Title(genAlbum.Name),
				Cover:    genAlbum.Cover,
			}

			albumID, err := g.db.CreateAlbum(ctx, album)
			if err != nil {
				log.Fatalf("failed to create album: %v", err)
			}

			log.Printf("Created album %s", strings.Title(genAlbum.Name))

			trackID := 1
			for i := 0; i < songs; i++ {
				var genSong *gen.Song
				gofakeit.Struct(&genSong)

				song := sqlc.CreateSongParams{
					AlbumID:  albumID,
					ArtistID: artistsID,
					Title:    strings.Title(genSong.Title),
					Length:   genSong.Length,
					Track:    sql.NullInt32{Int32: int32(trackID), Valid: true},
					Path:     genSong.Path,
					Mtime:    genSong.Mtime,
				}

				err := g.db.CreateSong(ctx, song)
				if err != nil {
					log.Fatalf("failed to create song: %v", err)
				}

				log.Printf("Created song %s", strings.Title(genSong.Title))
				trackID++
			}
		}
	}

	return nil
}
