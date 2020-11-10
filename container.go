package confucius

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"sync"
)

type (
	Service interface {
		Serve(handler *mux.Router) error
		Stop()
		Init() error
	}

	Middleware interface {
		Service
		Middleware(handlerFunc http.HandlerFunc) http.HandlerFunc
	}

	Container interface {
		Get(service string) (*Service, Status)
		Set(name string, service Service)
		Has(service string) bool
		Service
	}

	ServerContainer interface {
		Container
		SetMainService(name string, service Service)
	}

	serviceContainer struct {
		containerLock sync.Mutex
		services      []*containerEntry
		mainService   *containerEntry
		errorChannel  chan error
	}
)

func NewContainer() ServerContainer {
	return &serviceContainer{
		services:     make([]*containerEntry, 0),
		errorChannel: make(chan error),
	}
}

func (s *serviceContainer) SetMainService(name string, service Service) {
	s.containerLock.Lock()
	s.mainService = &containerEntry{
		name:    name,
		service: service,
		status:  Inactive,
	}
	s.containerLock.Unlock()
}

func (s *serviceContainer) Get(service string) (*Service, Status) {
	s.containerLock.Lock()
	defer s.containerLock.Unlock()

	if s.mainService.name == service {
		return &s.mainService.service, s.mainService.getStatus()
	}

	for _, e := range s.services {
		if e.name == service {
			return &e.service, e.getStatus()
		}
	}

	return nil, Undefined
}

func (s *serviceContainer) Set(name string, service Service) {
	s.containerLock.Lock()
	s.services = append(s.services, &containerEntry{
		name:    name,
		service: service,
		status:  Inactive,
	})
	s.containerLock.Unlock()
}

func (s *serviceContainer) Has(service string) bool {
	s.containerLock.Lock()
	defer s.containerLock.Unlock()

	for _, entry := range s.services {
		if entry.name == service {
			return true
		}
	}

	return false
}

func (s *serviceContainer) Serve(handler *mux.Router) error {
	running := 0
	for _, entry := range s.services {
		if entry.hasStatus(Ok) {
			running++
			go func(e *containerEntry) {
				e.setStatus(Serving)
				if err := e.service.Serve(handler); err != nil {
					s.errorChannel <- errors.Wrap(err, fmt.Sprintf("[%s]", e.name))
				}
				e.setStatus(Stopped)
			}(entry)
		}
	}

	if s.mainService.hasStatus(Ok) {
		running++
		s.mainService.setStatus(Serving)
		if err := s.mainService.service.Serve(handler); err != nil {
			s.errorChannel <- errors.Wrap(err, fmt.Sprintf("[%s]", s.mainService.name))
		}
		s.mainService.setStatus(Stopped)
	}

	// simple handler to handle empty configs
	if running == 0 {
		return nil
	}

	for fail := range s.errorChannel {
		s.Stop()
		return fail
	}

	return nil
}

func (s *serviceContainer) Stop() {
	for _, entry := range s.services {
		if entry.hasStatus(Serving) {
			entry.setStatus(Stopping)
			entry.service.Stop()
			entry.setStatus(Stopped)
		}
	}
}

func (s *serviceContainer) Init() error {
	if err := initService(s.mainService); err != nil {
		return err
	}

	for _, entry := range s.services {
		if err := initService(entry); err != nil {
			return err
		}
	}

	return nil
}

func initService(entry *containerEntry) error {
	if entry.getStatus() >= Ok {
		return fmt.Errorf("service [%s] has already been configured", entry.name)
	}

	err := entry.service.Init()

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("service [%s] cannot be initialized", entry.name))
	}

	entry.setStatus(Ok)

	return nil
}
