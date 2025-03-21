package entity

type Song struct {
	ID          int64  `json:"id" db:"id"`
	Group       string `json:"group" db:"group"`
	Title       string `json:"title" db:"title"`
	ReleaseDate string `json:"releaseDate" db:"release_date"`
	Text        string `json:"text" db:"text"`
	Link        string `json:"link" db:"link"`
}

type SongFilter struct {
	Group  string `json:"group"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type VersePagination struct {
	SongID int64 `json:"song_id"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}
