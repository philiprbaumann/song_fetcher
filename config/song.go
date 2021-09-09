package config

import (
	"bytes"
	"fmt"
)

type Song struct {
	Artist string `json:"artist"`
	Name   string `json:"name"`
}

func (s Song) String() string {
	return fmt.Sprintf("%v / %v\n", s.Artist, s.Name)
}

type ResponseSong struct {
	Released string `json:"released_at"`
	Songs    []Song `json:"songs"`
}

func (rs ResponseSong) String() string {
	var buff bytes.Buffer
	buff.WriteString(rs.Released + ":\n")
	for _, s := range rs.Songs {
		buff.WriteString(s.String() + "\n")
	}
	return fmt.Sprint(buff.String())
}

type Response struct {
	Data []ResponseSong
}

func (r Response) String() string {
	var buff bytes.Buffer
	for _, rs := range r.Data {
		buff.WriteString(rs.String() + "\n")
	}
	return fmt.Sprint(buff.String())
}
