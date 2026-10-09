package core

import (
	"bytes"
	"fmt"
	"strings"
)

const CRLF string = "\r\n"

var RespNil = []byte("$-1\r\n")

func readLen(data []byte) (int, int) {
	res, pos, _ := readInt64(data)
	return int(res), pos
}

// +<data>\r\n
func readSimpleString(data []byte) (string, int, error) {
	pos := 1
	for data[pos] != '\r' {
		pos++
	}
	return string(data[1:pos]), pos + 2, nil
}

// $11\r\nthangnguyen\r\n
func readBulkString(data []byte) (interface{}, int, error) {
	// null value - $-1\r\n -> minimum length of 5
	if len(data) < 6 {
		return "", len(data), ErrInvalidData
	}
	// empty string - $0\r\n\r\n
	if int(data[1]) == 0 {
		return "", len(data), nil
	}

	length, pos := readLen(data)
	if length == -1 {
		// null value
		return nil, len(data), nil
	}
	str := string(data[pos : pos+length])

	return str, pos + length + 2, nil
}

// :-123\r\n
func readInt64(data []byte) (int64, int, error) {
	if len(data) < 4 {
		return 0, 0, ErrInvalidData
	}
	var num int64 = 0
	var sign int64 = 1
	pos := 1
	if data[pos] == '-' {
		sign = -1
		pos++
	} else if data[pos] == '+' {
		pos++
	}
	for ; data[pos] != '\r'; pos++ {
		num = num*10 + int64(data[pos]-'0')
	}
	num = sign * num
	return num, pos + 2, nil
}

// *<len>\r\n<element 1><element 2>...
func readArray(data []byte) (interface{}, int, error) {
	if len(data) < 4 {
		return nil, 0, ErrInvalidData
	}

	// empty array
	if int(data[1]) == 0 {
		return make([]interface{}, 0), 0, nil
	}

	size, pos := readLen(data)

	if size == -1 {
		// null array
		return nil, len(data), nil
	}

	arr := make([]interface{}, size)
	for i := 0; i < size; i++ {
		// for each element of array, get its value from the data slice
		elem, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, len(data), err
		}
		arr[i] = elem
		pos += delta
	}

	return arr, pos, nil
}

// -<error>\r\n
func readError(data []byte) (interface{}, int, error) {
	if len(data) < 3 {
		return nil, 0, ErrInvalidData
	}
	// empty error message
	if len(data) == 3 {
		return "", 0, nil
	}

	pos := 1
	for data[pos] != '\r' {
		pos++
	}
	message := string(data[1:pos])
	return message, pos + 2, nil
}

// RESP format data -> raw data
func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, ErrNoData
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
	return nil, 0, nil
}

func Decode(data []byte) (interface{}, error) {
	res, _, err := DecodeOne(data)
	return res, err
}

// "hello" -> $5\r\nhello\r\n
func encodeString(s string) []byte {
	return fmt.Appendf(nil, "$%d%s%s%s", len(s), CRLF, s, CRLF)
}

func encodeStringArray(sa []string) []byte {
	var b []byte = make([]byte, 0)

	buf := bytes.NewBuffer(b)
	for _, s := range sa {
		buf.Write(encodeString(s))
	}
	return fmt.Appendf(nil, "*%d\r\n%s", len(sa), buf.Bytes())
}

// raw data -> RESP format data
func Encode(value interface{}, isSimpleString bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimpleString {
			return fmt.Appendf(nil, "+%s%s", v, CRLF)
		}
		// return []byte(fmt.Sprintf("$%d%s%s%s", len(v), CRLF, v, CRLF)) -> not optimized
		return encodeString(v)
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

func ParseCmd(data []byte) (*Command, error) {
	value, err := Decode(data)
	if err != nil {
		return nil, err
	}

	array := value.([]interface{})
	tokens := make([]string, len(array))
	for i := range tokens {
		tokens[i] = array[i].(string)
	}
	res := &Command{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}
	return res, nil
}
