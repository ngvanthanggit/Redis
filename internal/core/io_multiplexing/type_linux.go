//go:build linux

package io_multiplexing

import "syscall"

// convert application-level event to Linux Epoll Event
func (e Event) toNative() syscall.EpollEvent {
	return syscall.EpollEvent{}
}

// convert Linux Epoll Event to application-level event
func createEvent(ep syscall.EpollEvent) Event {
	return Event{}
}
