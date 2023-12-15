package pool

import (
	"loadbalancer/internal/service"
	"sync/atomic"
)

type Pool struct {
	Services []*service.Service
	Current  int64
}

func (s *Pool) nextInd() int {
	return int(atomic.AddInt64(&s.Current, 1)) % len(s.Services)
}

func (s *Pool) NextAlive() *service.Service {
	for i := 0; i < len(s.Services); i++ {
		n := s.nextInd()
		if s.Services[n].IsAlive() {
			return s.Services[n]
		}
	}

	return nil
}

func (s *Pool) AddService(sc *service.Service) {
	s.Services = append(s.Services, sc)
}
