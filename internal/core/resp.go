package core

import (
	"bytes"
	"errors"
	"fmt"
)

const CRLF string = "\r\n"

var RespNil = []byte("$-1\r\n")

func readSimpleString(data []byte) (interface{}, int, error) {
	return nil, 0, nil
}

func readBulkString(data []byte) (interface{}, int, error) {
	return nil, 0, nil
}

func readInt64(data []byte) (interface{}, int, error) {
	return nil, 0, nil
}

func readArray(data []byte) (interface{}, int, error) {
	return nil, 0, nil
}

func readError(data []byte) (interface{}, int, error) {
	return nil, 0, nil
}

// RESP format data -> raw data
func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}
	switch data[0] {
	case '+': // simple string
		return readSimpleString(data)
	case '$': // bulk string
		return readBulkString(data)
	case ':': // integer
		return readInt64(data)
	case '*': // array
		return readArray(data)
	case '-': // error
		return readError(data)
	}
}

func Decode(data []byte) (interface{}, error) {
	res, _, err := DecodeOne(data)
	return res, err
}

func encodeStringArray(sa []string) []byte {
	var b []byte = make([]byte, 0)

	for _, s := range sa {
		b = append(b, []byte(s)...)
	}
	return b
}

// raw data -> RESP format data
func Encode(value interface{}, isSimpleString bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimpleString {
			return fmt.Appendf(nil, "+%s%s", v, CRLF)
		}
		// return []byte(fmt.Sprintf("$%d%s%s%s", len(v), CRLF, v, CRLF)) -> not optimized
		return fmt.Appendf(nil, "$%d%s%s%s", len(v), CRLF, v, CRLF)
	case int64, int32, int16, int8, int:
		return fmt.Appendf(nil, ":%d%s", v, CRLF)
	case []string:
		return encodeStringArray(value.([]string))
	case [][]string:
		var b []byte
		// growing bytes storage - may replace with a simple byte slice (but i want to try :>)
		buf := bytes.NewBuffer(b)
		for _, sa := range v {
			buf.Write(encodeStringArray(sa))
		}
		return fmt.Appendf(nil, "*%d%s%s", len(v), CRLF, buf.Bytes())
	case []interface{}:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, x := range v {
			buf.Write(Encode(x, false))
		}
		return fmt.Appendf(nil, "*%d%s%s", len(v), CRLF, buf.Bytes())
	case error:
		return fmt.Appendf(nil, "-%s%s", v, CRLF)
	default:
		return RespNil
	}
}
