package repoerrs

import "errors"

var (
	ErrNotFound           = errors.New("song not found")
	ErrInsertFailed       = errors.New("failed to insert song")
	ErrUpdateFailed       = errors.New("failed to update song")
	ErrDeleteFailed       = errors.New("failed to delete song")
	ErrFetchSongsFailed   = errors.New("failed to fetch songs")
	ErrFetchVersesFailed  = errors.New("failed to fetch song verses")
	ErrStartTxFailed      = errors.New("failed to start transaction")
	ErrCommitTxFailed     = errors.New("failed to commit transaction")
	ErrRollbackTxFailed   = errors.New("failed to rollback transaction")
	ErrRowsAffectedFailed = errors.New("failed to check affected rows")
)
