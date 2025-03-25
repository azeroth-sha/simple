package grace

import (
	"github.com/azeroth-sha/simple"
	"sync"
	"sync/atomic"
	"time"
)

type Server interface {
	Start() error
	Stop() error
}

type service struct {
	running  int32
	name     string
	server   Server
	priority int
	logger   simple.Logger
	wait     *sync.WaitGroup
}

func (s *service) start() {
	go s.run()
}

func (s *service) run() {
	if atomic.SwapInt32(&s.running, 1) != 0 {
		return
	}
	s.wait.Add(1)
	defer s.wait.Done()
	defer s.stop()
	defer s.logger.Infof("service %s exited", s.name)
	s.logger.Infof("service %s starting", s.name)
	for atomic.LoadInt32(&s.running) == 1 {
		func() {
			defer func() {
				if rec := recover(); rec != nil {
					s.logger.Errorf("service %s panic: %v", s.name, rec)
				}
			}()
			if err := s.server.Start(); err != nil {
				s.logger.Errorf("service %s start error: %v", s.name, err)
				if atomic.LoadInt32(&s.running) == 1 {
					time.Sleep(time.Second)
				}
			}
		}()
	}
}

func (s *service) stop() {
	if atomic.SwapInt32(&s.running, 0) != 1 {
		return
	}
	s.logger.Infof("service %s stop", s.name)
	defer s.logger.Infof("service %s stopped", s.name)
	if err := s.server.Stop(); err != nil {
		s.logger.Errorf("service %s stop error: %v", s.name, err)
	}
}

func newServ(name string, server Server, logger simple.Logger, wait *sync.WaitGroup, priority ...int) *service {
	serv := &service{
		name:     name,
		server:   server,
		logger:   logger,
		wait:     wait,
		priority: 0,
	}
	if len(priority) > 0 {
		serv.priority = priority[0]
	}
	return serv
}
