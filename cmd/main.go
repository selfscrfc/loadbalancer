package main

import (
	"fmt"
	"loadbalancer/config"
	"loadbalancer/internal/balancer"
	"loadbalancer/internal/pool"
	"loadbalancer/internal/service"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err.Error())
	}

	balancer.Pl = &pool.Pool{}
	for _, e := range cfg.Services {
		url_, err := url.Parse(e)
		if err != nil {
			log.Fatal(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(url_)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
			log.Printf("Connection error, is service up?")
			for _, sc := range balancer.Pl.Services {
				if *sc.Url == *url_ {
					sc.SetAlive(false)
				}
			}
			balancer.Balance(w, r)
		}

		balancer.Pl.AddService(&service.Service{
			Url:    url_,
			Alive:  true,
			RProxy: proxy,
		})
		log.Printf("Service up %s\n", url_)
	}

	port, _ := strconv.Atoi(cfg.Port)
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(balancer.Balance),
	}

	go balancer.HealthCheck()
	log.Printf("Load balancer started at :%s", cfg.Port)
	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
