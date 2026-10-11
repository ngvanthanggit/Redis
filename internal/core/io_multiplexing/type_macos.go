//go:build darwin

package io_multiplexing

import "syscall"

func (e Event) toNative() syscall.Kevent_t {
	return syscall.Kevent_t{}
}

func createEvent(kq syscall.Kevent_t) Event {
	return Event{}
}
