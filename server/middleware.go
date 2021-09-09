package server

import (
	"github.com/valyala/fasthttp"
)

func (ws *Server) Recovery(next func(ctx *fasthttp.RequestCtx)) func(ctx *fasthttp.RequestCtx) {
	fn := func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if rvr := recover(); rvr != nil {
				ctx.Error("Internal Server Error", 500)
			}
		}()
		next(ctx)
	}
	return fn
}
