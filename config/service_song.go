package config

import (
	"fmt"
)

//ServiceSong struct for JSON unmarshalling.
type ServiceSong struct {
	ID         string `json:"song_id"`
	ReleasedAt string `json:"released_at"`
	Duration   string `json:"duration"`
	Artist     string `json:"artist"`
	Name       string `json:"name"`
	Stats      Stats  `json:"stats"`
}

func (song ServiceSong) String() string {
	return fmt.Sprintf("ServiceSong: %v / %v / %v / %v / %v", song.ID, song.ReleasedAt, song.Duration, song.Artist, song.Name)
}

//Stats object inside ServiceSong struct for JSON unmarshalling.
type Stats struct {
	LastPlayedAt int `json:"last_played_at"`
	TimesPlayed  int `json:"times_plaed"`
	GlobalRank   int `json:"global_rank"`
}
