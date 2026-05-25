package proxy

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

// Package provides information about each kindService(in K8S)
type KindService struct {
	URL          *url.URL
	Alive        bool
	ActiveConns  int64 // number of active connections at any point of time
	ReverseProxy *httputil.ReverseProxy

	// http handler(on a Goroutine) need access to these fields through mutex
	accessGate sync.RWMutex
}

func NewKindService(URL *url.URL) *KindService {
	return &KindService{
		URL:          URL,
		Alive:        true, // assume, kindService is alive initially
		ReverseProxy: httputil.NewSingleHostReverseProxy(URL),
	}
}

// SetAliveness() updates the backend status(for alivenesss)
func (b *KindService) SetAliveness(alivenessFlag bool) {
	b.accessGate.Lock()
	b.Alive = alivenessFlag
	b.accessGate.Unlock()
}

// IsAlive() checks whether backend services is alive
func (b *KindService) IsAlive() bool {
	b.accessGate.Lock()
	alivenessStatus := b.Alive
	b.accessGate.Unlock()
	return alivenessStatus
}
