package server

import (
	"net"

	"github.com/fasthttp/router"
	"github.com/oklog/run"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
)

type Server struct {
	Config  WebConfig
	Address string
	ln      net.Listener
	router  *router.Router
	debug   bool
}

func NewServer(config WebConfig) *Server {
	return &Server{
		Config:  config,
		Address: ServerAddr,
		router:  router.New(),
		debug:   true,
	}
}

func (s *Server) Close() {
	_ = s.ln.Close()
}

func (s *Server) Run() (err error) {
	s.muxRouter()

	s.ln, err = reuseport.Listen("tcp4", s.Address)
	log.Info("Starting server on " + s.Address)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	ws := &fasthttp.Server{
		Handler:            s.router.Handler,
		Name:               s.Config.Name,
		ReadBufferSize:     s.Config.ReadBufferSize,
		MaxConnsPerIP:      s.Config.MaxConnsPerIP,
		MaxRequestsPerConn: s.Config.MaxRequestsPerConn,
		MaxRequestBodySize: s.Config.MaxRequestBodySize,
		Concurrency:        s.Config.Concurrency,
	}

	var g run.Group
	g.Add(func() error {
		return ws.Serve(s.ln)
	}, func(e error) {
		_ = s.ln.Close()
	})
	return g.Run()
}
