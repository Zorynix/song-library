package app

import (
	"errors"
	"time"

	logger "github.com/Zorynix/song-library/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultAttempts = 20
	defaultTimeout  = time.Second
)

func RunMigrations(databaseURL string) error {
	if databaseURL == "" {
		return errors.New("database URL is empty")
	}

	var (
		attempts = defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", databaseURL)
		if err == nil {
			break
		}

		logger.Logger.Warn().Err(err).Int("attempts_left", attempts).Msg("Trying to connect to database for migrations")
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		return err
	}

	err = m.Up()
	defer func() { _, _ = m.Close() }()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logger.Logger.Info().Msg("No migration changes detected")
	} else {
		logger.Logger.Info().Msg("Migrations applied successfully")
	}

	return nil
}
