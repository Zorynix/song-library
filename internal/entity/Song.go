package entity

type Song struct {
	ID     int64  `json:"id"`
	Group  string `json:"group"`
	Title  string `json:"song"`
	Lyrics string `json:"lyrics,omitempty"`
}

type SongFilter struct {
	Group  string `json:"group,omitempty"`
	Title  string `json:"song,omitempty"`
	Lyrics string `json:"lyrics,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

type VersePagination struct {
	SongID int64 `json:"song_id"`
	Limit  int   `json:"limit,omitempty"`
	Offset int   `json:"offset,omitempty"`
}
