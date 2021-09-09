package server

import (
	"encoding/json"
	"time"

	config "github.com/prb-releases/config"

	"github.com/prb-releases/client"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func (s *Server) fetchSongHandler() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		// Parse parameters.
		artist := string(ctx.QueryArgs().Peek("artist"))
		start, err := time.Parse("2006-01-02", string(ctx.QueryArgs().Peek("from")))
		if err != nil {
			ctx.SetContentType("text/plain; charset=utf8")
			ctx.SetStatusCode(400)
			ctx.SetBody([]byte(`Invalid start time.`))
			return
		}
		end, err := time.Parse("2006-01-02", string(ctx.QueryArgs().Peek("until")))
		if err != nil {
			ctx.SetContentType("text/plain; charset=utf8")
			ctx.SetStatusCode(400)
			ctx.SetBody([]byte(`Invalid end time.`))
			return
		}
		if start.After(end) {
			ctx.SetContentType("text/plain; charset=utf8")
			ctx.SetStatusCode(400)
			ctx.SetBody([]byte(`Start time must be before end time.`))
			return
		}

		//Check for valid secret.
		c := client.New()
		if c.Secret == "XXX" {
			ctx.SetContentType("text/plain; charset=utf8")
			ctx.SetStatusCode(500)
			ctx.SetBody([]byte(`There was no secret set in the API. Please navigate to client/client.go and set the secret.`))
			return
		}

		// Fetch songs from release provider.
		log.Info("Fetching songs from release provider...")
		providerResp := c.GetMultipleSongs(start, end, &artist)
		if len(providerResp) == 0 {
			ctx.SetStatusCode(404)
			ctx.SetBody([]byte(`No songs found.`))
			return
		}

		var apiResp []*config.ResponseSong
		for k, v := range providerResp {
			rs := config.ResponseSong{Released: k, Songs: v}
			apiResp = append(apiResp, &rs)
		}
		response, err := json.Marshal(apiResp)
		if err != nil {
			log.Error("Failed to unmarshal API response.")
		}
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(200)
		ctx.SetBody([]byte(response))
	}
}
