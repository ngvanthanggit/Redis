package core

import (
	"fmt"
	"syscall"
)

func cmdPING(args []string) []byte {
	var res []byte
	if len(args) > 1 {
		return Encode(ErrNumberOfArgs, false)
	}

	if len(args) == 0 {
		res = Encode("PONG", true)
	} else {
		res = Encode(args[0], false)
	}
	return res
}

func ExecuteAndResponse(cmd *Command, connFd int) error {
	var res []byte

	switch cmd.Cmd {
	case "PING":
		res = cmdPING(cmd.Args)
	default:
		res = Encode(fmt.Errorf("ERR unknown command '%s'", cmd.Cmd), false)
	}
	_, err := syscall.Write(connFd, res)
	return err
}
