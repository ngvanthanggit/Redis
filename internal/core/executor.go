package core

import (
	"fmt"
	"ngvanthanggit/internal/constant"
	"strconv"
	"syscall"
	"time"
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

func cmdSET(args []string) []byte {
	if len(args) < 2 {
		return Encode(ErrNumberOfArgs, false)
	}

	// process the arguments
	// SET key value EX ttl
	var key, value string
	var ttlMs int64
	key, value = args[0], args[1]
	ttlMs = -1 // default if not specified

	if len(args) > 2 {
		if len(args) != 4 {
			return Encode(ErrNumberOfArgs, false)
		}
		ttlSec, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return Encode(err, false)
		}

		ttlMs = ttlSec * 1000
	}

	dictStore.Set(key, dictStore.NewObj(key, value, ttlMs))

	return constant.RespOk
}

func cmdGET(args []string) []byte {
	if len(args) < 1 {
		return Encode(ErrNumberOfArgs, false)
	}

	key := args[0]
	value := dictStore.Get(key)
	if value == nil {
		return constant.RespNil
	}

	return Encode(value.Value, false)
}

func cmdTTL(args []string) []byte {
	if len(args) < 1 {
		return Encode(ErrNumberOfArgs, false)
	}

	key := args[0]
	if dictStore.HasExpired(key) {
		return constant.TtlKeyNotExist
	}

	exp, exist := dictStore.GetExpiry(key)
	if !exist {
		return constant.TtlKeyExistNoExpire
	}

	ttlMs := exp - uint64(time.Now().UnixMilli())
	ttlSec := uint64(ttlMs / 1000)

	return Encode(ttlSec, false)
}

func ExecuteAndResponse(cmd *Command, connFd int) error {
	var res []byte

	switch cmd.Cmd {
	// normal commands
	case "PING":
		res = cmdPING(cmd.Args)
	case "SET":
		res = cmdSET(cmd.Args)
	case "GET":
		res = cmdGET(cmd.Args)
	case "TTL":
		res = cmdTTL(cmd.Args)
	// simple set commands
	case "SADD":
		res = cmdSADD(cmd.Args)
	case "SREM":
		res = cmdSREM(cmd.Args)
	case "SISMEMBER":
		res = cmdSISMEMBER(cmd.Args)
	case "SMEMBERS":
		res = cmdSMEMBERS(cmd.Args)
	// sorted set commands
	case "ZADD":
		res = cmdZADD(cmd.Args)
	case "ZSCORE":
		res = cmdZSCORE(cmd.Args)
	case "ZRANK":
		res = cmdZRANK(cmd.Args)
	// count-min sketch commands
	case "CMS.INITBYDIM":
		res = cmdCMSINITBYDIM(cmd.Args)
	case "CMS.INITBYPROB":
		res = cmdCMSINITBYPROB(cmd.Args)
	case "CMS.INCRBY":
		res = cmdCMSINCRBY(cmd.Args)
	case "CMS.QUERY":
		res = cmdCMSQUERY(cmd.Args)
	// bloom filter
	case "BF.RESERVE":
		res = cmdBFRESERVE(cmd.Args)
	case "BF.MADD":
		res = cmdBFMADD(cmd.Args)
	case "BF.EXISTS":
		res = cmdBFEXISTS(cmd.Args)
	default:
		res = Encode(fmt.Errorf("ERR unknown command '%s'", cmd.Cmd), false)
	}
	_, err := syscall.Write(connFd, res)
	return err
}
