package balancer

import (
	"loadbalancer/internal/pool"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

var Pl *pool.Pool

func Balance(w http.ResponseWriter, r *http.Request) {
	n := Pl.NextAlive()
	if n != nil {
		n.RProxy.ServeHTTP(w, r)
		return
	}

	http.Error(w, "No services available", http.StatusServiceUnavailable)
}

func HealthCheck() {
	for {
		for _, s := range Pl.Services {
			status := "up"
			alive := isBackendAlive(s.Url)
			s.SetAlive(alive)
			if !alive {
				status = "down"
			}
			log.Printf("%s [%s]\n", s.Url, status)
		}
		time.Sleep(20 * time.Second)
	}

}

func isBackendAlive(u *url.URL) bool {
	conn, err := net.DialTimeout("tcp", u.Host, 2*time.Second)
	if err != nil {
		log.Println("Ping error: ", err)
		return false
	}
	_ = conn.Close()
	return true
}
