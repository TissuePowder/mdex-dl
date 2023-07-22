package main

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	// "golang.org/x/time/rate"
	rate "github.com/projectdiscovery/ratelimit"
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
			KeepAlive: 30 * time.Second,
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
		"/at-home/server/": rate.New(context.Background(), 38, time.Minute),
	}

	return &MdClient{
		Client:         client,
		RateLimiterMap: rlMap,
	}

}

func (c *MdClient) Get(url string) (*http.Response, error) {

	if strings.Contains(url, "/at-home/server/") {
		// fmt.Println("We are rate limiting!")
		// c.Mutex.Lock()
		limiter := c.RateLimiterMap["/at-home/server/"]
		// c.Mutex.Unlock()
		limiter.Take()
	}

	res, err := c.Client.Get(url)

	// fmt.Println(res.StatusCode)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	return res, err

}
