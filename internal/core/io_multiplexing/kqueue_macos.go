//go:build darwin

package io_multiplexing

import (
	"log"
	"ngvanthanggit/internal/config"
	"syscall"
)

type KQueue struct {
	fd            int
	kqEvents      []syscall.Kevent_t
	genericEvents []Event
}

func CreateIOMultiplexer() (*KQueue, error) {
	kqueueFD, err := syscall.Kqueue()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &KQueue{
		fd:            kqueueFD,
		kqEvents:      make([]syscall.Kevent_t, config.MaxConnection),
		genericEvents: make([]Event, config.MaxConnection),
	}, nil
}

func (kq *KQueue) Monitor(event Event) error {
	return nil
}

func (kq *KQueue) Wait() ([]Event, error) {
	return nil, nil
}

func (kq *KQueue) Close() error {
	return nil
}
