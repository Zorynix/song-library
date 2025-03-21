package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Zorynix/song-library/internal/entity"
	"github.com/Zorynix/song-library/internal/repo/pgdb"
)

type SongRepo interface {
	GetSongs(ctx context.Context, filter entity.SongFilter) ([]entity.Song, error)
	GetSongVerses(ctx context.Context, pagination entity.VersePagination) ([]string, error)
	DeleteSong(ctx context.Context, id int64) error
	UpdateSong(ctx context.Context, song entity.Song) error
	AddSong(ctx context.Context, song entity.Song) (entity.Song, error)
}

type Repositories struct {
	Song SongRepo
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Song: pgdb.NewSongRepo(db),
	}
}
