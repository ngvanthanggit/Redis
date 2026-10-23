package server

import (
	"io"
	"log"
	"net"
	"ngvanthanggit/internal/config"
	"ngvanthanggit/internal/constant"
	"ngvanthanggit/internal/core"
	"ngvanthanggit/internal/core/io_multiplexing"
	"syscall"
	"time"
)

func readCommand(fd int) (*core.Command, error) {
	var buf = make([]byte, 512)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return nil, err
	}
	// 0 bytes were read
	if n == 0 {
		return nil, io.EOF
	}
	// parse the command here using the resp.ParseCmd(buf)
	cmd, err := core.ParseCmd(buf)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func respond(data string, fd int) error {
	if _, err := syscall.Write(fd, []byte(data)); err != nil {
		return err
	}
	return nil
}

func RunIoMultiplexingServer() {
	log.Println("starting an I/O Multiplexing TCP server on", config.Port)

	listener, err := net.Listen(config.Protocol, config.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// Get the file descriptor of the Listener
	// 1. Use TCP listener
	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		log.Fatal("listener is not a TCP Listener")
	}
	// 2. Get the file descriptor of the listener
	listenerFile, err := tcpListener.File()
	if err != nil {
		log.Fatal(err)
	}
	defer listenerFile.Close()

	serverFd := int(listenerFile.Fd())

	// create an ioMultiplexer instance
	ioMultiplexer, err := io_multiplexing.CreateIOMultiplexer()
	if err != nil {
		log.Fatal(err)
	}
	defer ioMultiplexer.Close()

	// Monitor 'read' events from the Server FD
	if err = ioMultiplexer.Monitor(io_multiplexing.Event{
		Fd: serverFd,
		Op: io_multiplexing.OpRead,
	}); err != nil {
		log.Fatal(err)
	}

	lastActiveExpireExecTime := time.Now()

	var events = make([]io_multiplexing.Event, config.MaxConnection)
	for {
		if time.Now().After(lastActiveExpireExecTime.Add(constant.ActiveExpireFrequency)) {
			core.ActiveDeleteExpiredKeys()
			lastActiveExpireExecTime = time.Now()
		}
		// wait for at least one FD in the monitoring list to be ready for I/O
		// this is a blocking call
		events, err = ioMultiplexer.Wait()
		if err != nil {
			continue
		}

		for i := range len(events) {
			if events[i].Fd == serverFd {
				log.Println("new client is trying to connect")

				// create a new connection
				// the syscall.Accept() return the fd of the new socket
				connFd, _, err := syscall.Accept(serverFd)
				if err != nil {
					log.Fatal(err)
					continue
				}

				log.Println("set up a new connection")
				// let the ioMultiplexer monitor Read events for Server FD
				if err = ioMultiplexer.Monitor(
					io_multiplexing.Event{
						Fd: connFd,
						Op: io_multiplexing.OpRead,
					},
				); err != nil {
					log.Fatal(err)
				}
			} else {
				// read the event command and execute
				cmd, err := readCommand(events[i].Fd)
				if err != nil {
					if err == io.EOF || err == syscall.ECONNRESET {
						log.Println("client disconnected")
						_ = syscall.Close(events[i].Fd)
						continue
					}
					log.Fatal(err)
					continue
				}

				// if err = respond(cmd, events[i].Fd); err != nil {
				// 	log.Fatal(err)
				// }

				// instead, we can use the executor.ExecuteAndResponse(cmd, events[i].Fd) here
				if err := core.ExecuteAndResponse(cmd, events[i].Fd); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
