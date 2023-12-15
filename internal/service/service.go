package service

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Service struct {
	Url    *url.URL
	Alive  bool
	mu     sync.Mutex
	RProxy *httputil.ReverseProxy
}

func (sc *Service) SetAlive(alive bool) {
	sc.mu.Lock()
	sc.Alive = alive
	sc.mu.Unlock()
}

func (sc *Service) IsAlive() bool {
	sc.mu.Lock()
	a := sc.Alive
	sc.mu.Unlock()
	return a
}
