package main

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

type MdClient struct {
	Client         *http.Client
	RateLimiterMap map[string]*rate.Limiter
	// Mutex          sync.Mutex
}

func NewMdClient() *MdClient {

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxConnsPerHost:       100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
	}

	rlMap := map[string]*rate.Limiter{
		"/at-home/server/": rate.NewLimiter(rate.Every(time.Minute/38), 1),
	}

	return &MdClient{
		Client:         client,
		RateLimiterMap: rlMap,
	}

}

func (c *MdClient) Get(url string) (*http.Response, error) {

	if strings.Contains(url, "/at-home/server/") {
		limiter := c.RateLimiterMap["/at-home/server/"]
		err := limiter.Wait(context.Background())
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept-Encoding", "identity")

	res, err := c.Client.Do(req)

	return res, err

}
