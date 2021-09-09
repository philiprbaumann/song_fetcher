package server

func (s *Server) muxRouter() {
	s.router.GET("/releases", s.Recovery(s.fetchSongHandler()))
}
