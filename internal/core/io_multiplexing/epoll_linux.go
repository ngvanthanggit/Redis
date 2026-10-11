//go:build linux

package io_multiplexing

import (
	"log"
	"ngvanthanggit/internal/config"
	"syscall"
)

type Epoll struct {
	fd            int
	epollEvents   []syscall.EpollEvent
	genericEvents []Event
}

func CreateIOMultiplexer() (*Epoll, error) {
	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Epoll{
		fd:            epollFD,
		epollEvents:   make([]syscall.EpollEvent, config.MaxConnection),
		genericEvents: make([]Event, config.MaxConnection),
	}, nil
}

func (ep *Epoll) Monitor(event Event) error {
	return nil
}

func (ep *Epoll) Wait() ([]Event, error) {
	return nil, nil
}

func (ep *Epoll) Close() error {
	return nil
}
