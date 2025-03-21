package pgdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/Zorynix/song-library/internal/entity"
	logger "github.com/Zorynix/song-library/internal/logger"
	repoerrs "github.com/Zorynix/song-library/internal/repo/repo_errors"
	"github.com/jmoiron/sqlx"
)

type SongRepo struct {
	db *sqlx.DB
}

func NewSongRepo(db *sqlx.DB) *SongRepo {
	return &SongRepo{db: db}
}

func (r *SongRepo) GetSongs(ctx context.Context, filter entity.SongFilter) ([]entity.Song, error) {
	logger.Logger.Debug().
		Str("group", filter.Group).
		Str("title", filter.Title).
		Str("text", filter.Text).
		Int("limit", filter.Limit).
		Int("offset", filter.Offset).
		Msg("Fetching songs with filter")

	var songs []entity.Song
	query := `SELECT id, "group", title, release_date, text, link FROM library.songs WHERE 1=1`
	var args []interface{}
	argIndex := 1

	if filter.Group != "" {
		query += fmt.Sprintf(" AND \"group\" ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Group+"%")
		argIndex++
	}
	if filter.Title != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Title+"%")
		argIndex++
	}
	if filter.Text != "" {
		query += fmt.Sprintf(" AND text ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Text+"%")
		argIndex++
	}

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	err := r.db.SelectContext(ctx, &songs, query, args...)
	if err != nil {
		logger.Logger.Error().Err(err).Msg(repoerrs.ErrFetchSongsFailed.Error())
		return nil, fmt.Errorf("%w: %v", repoerrs.ErrFetchSongsFailed, err)
	}

	logger.Logger.Info().Int("count", len(songs)).Msg("Songs fetched successfully")
	return songs, nil
}

func (r *SongRepo) GetSongVerses(ctx context.Context, pagination entity.VersePagination) ([]string, error) {
	logger.Logger.Debug().
		Int64("song_id", pagination.SongID).
		Int("limit", pagination.Limit).
		Int("offset", pagination.Offset).
		Msg("Fetching song verses")

	var text string
	err := r.db.GetContext(ctx, &text, `SELECT text FROM library.songs WHERE id = $1`, pagination.SongID)
	if err != nil {
		logger.Logger.Error().Err(err).Int64("song_id", pagination.SongID).Msg(repoerrs.ErrFetchVersesFailed.Error())
		return nil, fmt.Errorf("%w: %v", repoerrs.ErrFetchVersesFailed, err)
	}

	if text == "" {
		logger.Logger.Info().Int64("song_id", pagination.SongID).Msg("Song text is empty")
		return []string{}, nil
	}

	verses := strings.Split(text, "\n\n")
	if len(verses) == 0 {
		verses = strings.Split(text, "\n")
	}

	start := pagination.Offset
	if start >= len(verses) {
		logger.Logger.Warn().Int64("song_id", pagination.SongID).Int("offset", start).Msg("Offset exceeds verses count")
		return []string{}, nil
	}

	end := start + pagination.Limit
	if pagination.Limit <= 0 || end > len(verses) {
		end = len(verses)
	}

	logger.Logger.Info().
		Int64("song_id", pagination.SongID).
		Int("verse_count", len(verses[start:end])).
		Msg("Song verses fetched successfully")
	return verses[start:end], nil
}

func (r *SongRepo) DeleteSong(ctx context.Context, id int64) error {
	logger.Logger.Debug().Int64("id", id).Msg("Deleting song")

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Logger.Error().Err(err).Msg(repoerrs.ErrStartTxFailed.Error())
		return fmt.Errorf("%w: %v", repoerrs.ErrStartTxFailed, err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logger.Logger.Error().Err(rollbackErr).Msg(repoerrs.ErrRollbackTxFailed.Error())
				err = fmt.Errorf("%w: %v", repoerrs.ErrRollbackTxFailed, rollbackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			logger.Logger.Error().Err(commitErr).Msg(repoerrs.ErrCommitTxFailed.Error())
			err = fmt.Errorf("%w: %v", repoerrs.ErrCommitTxFailed, commitErr)
		}
	}()

	result, err := tx.ExecContext(ctx, `DELETE FROM library.songs WHERE id = $1`, id)
	if err != nil {
		logger.Logger.Error().Err(err).Int64("id", id).Msg(repoerrs.ErrDeleteFailed.Error())
		return fmt.Errorf("%w: %v", repoerrs.ErrDeleteFailed, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Logger.Error().Err(err).Int64("id", id).Msg(repoerrs.ErrRowsAffectedFailed.Error())
		return fmt.Errorf("%w: %v", repoerrs.ErrRowsAffectedFailed, err)
	}
	if rows == 0 {
		logger.Logger.Warn().Int64("id", id).Msg(repoerrs.ErrNotFound.Error())
		return repoerrs.ErrNotFound
	}

	logger.Logger.Info().Int64("id", id).Msg("Song deleted successfully")
	return nil
}

func (r *SongRepo) UpdateSong(ctx context.Context, song entity.Song) error {
	logger.Logger.Debug().
		Int64("id", song.ID).
		Str("group", song.Group).
		Str("title", song.Title).
		Msg("Updating song")

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Logger.Error().Err(err).Msg(repoerrs.ErrStartTxFailed.Error())
		return fmt.Errorf("%w: %v", repoerrs.ErrStartTxFailed, err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logger.Logger.Error().Err(rollbackErr).Msg(repoerrs.ErrRollbackTxFailed.Error())
				err = fmt.Errorf("%w: %v", repoerrs.ErrRollbackTxFailed, rollbackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			logger.Logger.Error().Err(commitErr).Msg(repoerrs.ErrCommitTxFailed.Error())
			err = fmt.Errorf("%w: %v", repoerrs.ErrCommitTxFailed, commitErr)
		}
	}()

	query := `
		UPDATE library.songs 
		SET "group" = $1, title = $2, release_date = $3, text = $4, link = $5 
		WHERE id = $6`
	result, err := tx.ExecContext(ctx, query, song.Group, song.Title, song.ReleaseDate, song.Text, song.Link, song.ID)
	if err != nil {
		logger.Logger.Error().Err(err).Int64("id", song.ID).Msg(repoerrs.ErrUpdateFailed.Error())
		return fmt.Errorf("%w: %v", repoerrs.ErrUpdateFailed, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Logger.Error().Err(err).Int64("id", song.ID).Msg(repoerrs.ErrRowsAffectedFailed.Error())
		return fmt.Errorf("%w: %v", repoerrs.ErrRowsAffectedFailed, err)
	}
	if rows == 0 {
		logger.Logger.Warn().Int64("id", song.ID).Msg(repoerrs.ErrNotFound.Error())
		return repoerrs.ErrNotFound
	}

	logger.Logger.Info().Int64("id", song.ID).Msg("Song updated successfully")
	return nil
}

func (r *SongRepo) AddSong(ctx context.Context, song entity.Song) (entity.Song, error) {
	logger.Logger.Debug().
		Str("group", song.Group).
		Str("title", song.Title).
		Msg("Adding new song")

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Logger.Error().Err(err).Msg(repoerrs.ErrStartTxFailed.Error())
		return entity.Song{}, fmt.Errorf("%w: %v", repoerrs.ErrStartTxFailed, err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logger.Logger.Error().Err(rollbackErr).Msg(repoerrs.ErrRollbackTxFailed.Error())
				err = fmt.Errorf("%w: %v", repoerrs.ErrRollbackTxFailed, rollbackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			logger.Logger.Error().Err(commitErr).Msg(repoerrs.ErrCommitTxFailed.Error())
			err = fmt.Errorf("%w: %v", repoerrs.ErrCommitTxFailed, commitErr)
		}
	}()

	query := `
		INSERT INTO library.songs ("group", title, release_date, text, link) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, "group", title, release_date, text, link`
	var createdSong entity.Song
	err = tx.GetContext(ctx, &createdSong, query, song.Group, song.Title, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		logger.Logger.Error().Err(err).Str("group", song.Group).Str("title", song.Title).Msg(repoerrs.ErrInsertFailed.Error())
		return entity.Song{}, fmt.Errorf("%w: %v", repoerrs.ErrInsertFailed, err)
	}

	logger.Logger.Info().Int64("id", createdSong.ID).Msg("Song added successfully")
	return createdSong, nil
}
