package services

import (
	"context"

	"github.com/Zorynix/song-library/internal/entity"
	"github.com/Zorynix/song-library/internal/repo"
)

type SongService interface {
	GetSongs(ctx context.Context, filter entity.SongFilter) ([]entity.Song, error)
	GetSongVerses(ctx context.Context, pagination entity.VersePagination) ([]string, error)
	DeleteSong(ctx context.Context, id int64) error
	UpdateSong(ctx context.Context, song entity.Song) error
	AddSong(ctx context.Context, song entity.Song) (entity.Song, error)
}

type Services struct {
	Song SongService
}

type ServicesDependencies struct {
	Repos *repo.Repositories
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Song: NewSongService(deps.Repos),
	}
}
