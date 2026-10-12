//go:build linux

package io_multiplexing

import "syscall"

// convert application-level event to Linux Epoll Event
func (e Event) toNative() syscall.EpollEvent {
	var event uint32
	switch e.Op {
	case OpWrite:
		event = syscall.EPOLLOUT
	case OpRead:
		event = syscall.EPOLLIN
	}
	return syscall.EpollEvent{
		Fd:     int32(e.Fd),
		Events: event,
	}
}

// convert Linux Epoll Event to application-level event
func createEvent(ep syscall.EpollEvent) Event {
	var op Operation
	switch ep.Events {
	case syscall.EPOLLOUT:
		op = OpWrite
	case syscall.EPOLLIN:
		op = OpRead
	}

	return Event{
		Fd: int(ep.Fd),
		Op: op,
	}
}
