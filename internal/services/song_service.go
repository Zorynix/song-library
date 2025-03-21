package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/Zorynix/song-library/internal/entity"
	errs "github.com/Zorynix/song-library/internal/errors"
	logger "github.com/Zorynix/song-library/internal/logger"
	"github.com/Zorynix/song-library/internal/repo"
	repoerrs "github.com/Zorynix/song-library/internal/repo/repo_errors"
)

type songService struct {
	repos       *repo.Repositories
	httpClient  *http.Client
	musicAPIURL string
}

func NewSongService(repos *repo.Repositories, musicAPIURL string) SongService {
	return &songService{
		repos:       repos,
		httpClient:  &http.Client{},
		musicAPIURL: musicAPIURL,
	}
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func (s *songService) GetSongs(ctx context.Context, filter entity.SongFilter) ([]entity.Song, error) {
	logger.Logger.Debug().
		Str("group", filter.Group).
		Str("title", filter.Title).
		Str("text", filter.Text).
		Int("limit", filter.Limit).
		Int("offset", filter.Offset).
		Msg("Fetching songs")

	songs, err := s.repos.Song.GetSongs(ctx, filter)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to fetch songs in service")
		return nil, errs.ErrInternal
	}

	logger.Logger.Info().Int("count", len(songs)).Msg("Songs fetched successfully in service")
	return songs, nil
}

func (s *songService) GetSongVerses(ctx context.Context, pagination entity.VersePagination) ([]string, error) {
	logger.Logger.Debug().
		Int64("song_id", pagination.SongID).
		Int("limit", pagination.Limit).
		Int("offset", pagination.Offset).
		Msg("Fetching song verses")

	if pagination.SongID <= 0 {
		logger.Logger.Error().Int64("song_id", pagination.SongID).Msg("Invalid song ID in service")
		return nil, errs.ErrInvalidInput
	}

	verses, err := s.repos.Song.GetSongVerses(ctx, pagination)
	if err != nil {
		logger.Logger.Error().Err(err).Int64("song_id", pagination.SongID).Msg("Failed to fetch song verses in service")
		return nil, errs.ErrInternal
	}

	logger.Logger.Info().
		Int64("song_id", pagination.SongID).
		Int("verse_count", len(verses)).
		Msg("Song verses fetched successfully in service")
	return verses, nil
}

func (s *songService) DeleteSong(ctx context.Context, id int64) error {
	logger.Logger.Debug().Int64("id", id).Msg("Deleting song")

	if id <= 0 {
		logger.Logger.Error().Int64("id", id).Msg("Invalid song ID in service")
		return errs.ErrInvalidInput
	}

	err := s.repos.Song.DeleteSong(ctx, id)
	if err != nil {
		logger.Logger.Error().Err(err).Int64("id", id).Msg("Failed to delete song in service")
		if errors.Is(err, repoerrs.ErrNotFound) {
			return errs.ErrNotFound
		}
		return errs.ErrInternal
	}

	logger.Logger.Info().Int64("id", id).Msg("Song deleted successfully in service")
	return nil
}

func (s *songService) UpdateSong(ctx context.Context, song entity.Song) error {
	logger.Logger.Debug().
		Int64("id", song.ID).
		Str("group", song.Group).
		Str("title", song.Title).
		Msg("Updating song")

	if song.ID <= 0 {
		logger.Logger.Error().Int64("id", song.ID).Msg("Invalid song ID in service")
		return errs.ErrInvalidInput
	}

	err := s.repos.Song.UpdateSong(ctx, song)
	if err != nil {
		logger.Logger.Error().Err(err).Int64("id", song.ID).Msg("Failed to update song in service")
		if errors.Is(err, repoerrs.ErrNotFound) {
			return errs.ErrNotFound
		}
		return errs.ErrInternal
	}

	logger.Logger.Info().Int64("id", song.ID).Msg("Song updated successfully in service")
	return nil
}

func (s *songService) AddSong(ctx context.Context, song entity.Song) (entity.Song, error) {
	logger.Logger.Debug().
		Str("group", song.Group).
		Str("title", song.Title).
		Msg("Adding new song")

	params := url.Values{}
	params.Add("group", song.Group)
	params.Add("song", song.Title)
	reqURL := s.musicAPIURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		logger.Logger.Error().Err(err).Str("url", reqURL).Msg("Failed to create request to music API")
		return entity.Song{}, errs.ErrInternal
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		logger.Logger.Error().Err(err).Str("url", reqURL).Msg("Failed to fetch data from music API")
		return entity.Song{}, errs.ErrInternal
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Logger.Error().
			Int("status", resp.StatusCode).
			Str("url", reqURL).
			Msg("Music API returned non-200 status")
		return entity.Song{}, errs.ErrInternal
	}

	var songDetail SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to decode music API response")
		return entity.Song{}, errs.ErrInternal
	}

	song.ReleaseDate = songDetail.ReleaseDate
	song.Text = songDetail.Text
	song.Link = songDetail.Link

	createdSong, err := s.repos.Song.AddSong(ctx, song)
	if err != nil {
		logger.Logger.Error().Err(err).Str("group", song.Group).Str("title", song.Title).Msg("Failed to add song in service")
		return entity.Song{}, errs.ErrInternal
	}

	logger.Logger.Info().Int64("id", createdSong.ID).Msg("Song added successfully in service")
	return createdSong, nil
}
