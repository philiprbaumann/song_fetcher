package client

import (
	"encoding/json"
	"sync"
	"time"

	config "github.com/prb-releases/config"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Client struct {
	BaseURI string
	Secret  string
}

func New() *Client {
	return &Client{
		BaseURI: "SERVER_NAME_HERE",
		Secret:  "XXX",
	}
}

func (w *Client) GetSongs(ch chan *fasthttp.Response, wg *sync.WaitGroup, url string) error {
	defer wg.Done()
	log.Info("Client is fetching songs at " + url)
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.SetRequestURI(url)
	req.Header.SetContentType("application/json")
	req.Header.SetMethod("GET")

	err := fasthttp.Do(req, resp)

	if err != nil {
		log.Error("POST REQUEST ERROR", err.Error())
		return err
	}

	log.Info("Fetching completed at " + url)
	out := fasthttp.AcquireResponse()
	resp.CopyTo(out)
	ch <- out
	return nil
}

func (c *Client) GetMultipleSongs(start time.Time, end time.Time, artist *string) map[string][]config.Song {
	// Execution time for total fetch.
	begin := time.Now()
	defer func() {
		log.Infof("Execution Time: %v\n\n", time.Since(begin))
	}()

	ch := make(chan *fasthttp.Response)
	var wg sync.WaitGroup
	var d time.Time

	// Add day queries until first of month is hit or end.
	for d = start; !d.After(end) && d.Day() != 1; d = d.AddDate(0, 0, 1) {
		wg.Add(1)
		go c.GetSongs(ch, &wg, "SERVER_NAME_HERE"+"/v1/songs/daily?released_at="+d.Format("2006-01-02")+"&api_key="+c.Secret)
	}
	// Add month queries until we exceed one month from the end.
	for ; !d.After(end.AddDate(0, -1, 0)); d = d.AddDate(0, 1, 0) {
		wg.Add(1)
		go c.GetSongs(ch, &wg, "SERVER_NAME_HERE"+"/v1/songs/monthly?released_at="+d.Format("2006-01")+"&api_key="+c.Secret)
	}
	// Add day queries for the remaining time.
	for ; !d.After(end); d = d.AddDate(0, 0, 1) {
		wg.Add(1)
		go c.GetSongs(ch, &wg, "SERVER_NAME_HERE"+"/v1/songs/daily?released_at="+d.Format("2006-01-02")+"&api_key="+c.Secret)
	}

	// Close channel in the background.
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Read from the channel as they come and until its closed.
	rp := make(map[string][]config.Song)
	for res := range ch {
		var ssArr []config.ServiceSong
		json.Unmarshal(res.Body(), &ssArr)
		for _, ss := range ssArr {
			if *artist != "" {
				if *artist == ss.Artist {
					s := config.Song{Artist: ss.Artist, Name: ss.Name}
					rp[ss.ReleasedAt] = append(rp[ss.ReleasedAt], s)
					log.Info(ss.String())
				}
			} else {
				s := config.Song{Artist: ss.Artist, Name: ss.Name}
				rp[ss.ReleasedAt] = append(rp[ss.ReleasedAt], s)
				log.Info(ss.String())
			}
		}
	}
	return rp
}
