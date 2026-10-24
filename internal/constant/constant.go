package constant

import "time"

var RespNil = []byte("$-1\r\n")
var RespOk = []byte("+OK\r\n")
var TtlKeyNotExist = []byte(":-2\r\n")
var TtlKeyExistNoExpire = []byte(":-1\r\n")
var ActiveExpireSampleSize = 20
var ActiveExpireThreshold = 0.1
var ActiveExpireFrequency = 100 * time.Millisecond
var BPlusTreeDegress = 4
