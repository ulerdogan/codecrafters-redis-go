package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Type byte

const (
	SimpleString Type = '+'
	BulkString   Type = '$'
	Array        Type = '*'
)

type Value struct {
	typ     Type
	command string
	args    []Value
	bytes   []byte
}

func (v Value) String() string {
	if v.typ == BulkString || v.typ == SimpleString {
		return string(v.bytes)
	}
	return ""
}

func (v Value) Array() []Value {
	if v.typ == Array {
		return v.args
	}

	return []Value{}
}

func DecodeRESP(byteStream *bufio.Reader) (*Value, error) {
	dataTypeByte, err := byteStream.ReadByte()
	if err != nil {
		return nil, err
	}

	switch string(dataTypeByte) {
	case "+":
		return decodeSimpleString(byteStream)
	case "$":
		return decodeBulkString(byteStream)
	case "*":
		return decodeArray(byteStream)
	}

	return nil, fmt.Errorf("invalid RESP data type byte: %s", string(dataTypeByte))
}

func decodeSimpleString(byteStream *bufio.Reader) (*Value, error) {
	readBytes, err := readUntilCRLF(byteStream)
	if err != nil {
		return nil, err
	}

	return &Value{
		typ:   SimpleString,
		bytes: readBytes,
	}, nil
}

func readUntilCRLF(byteStream *bufio.Reader) ([]byte, error) {
	readBytes := []byte{}

	for {
		b, err := byteStream.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		readBytes = append(readBytes, b...)
		if len(readBytes) >= 2 && readBytes[len(readBytes)-2] == '\r' {
			break
		}
	}

	return readBytes[:len(readBytes)-2], nil
}

func decodeBulkString(byteStream *bufio.Reader) (*Value, error) {
	readBytesForCount, err := readUntilCRLF(byteStream)
	if err != nil {
		return nil, fmt.Errorf("failed to read bulk string length: %s", err)
	}

	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return nil, fmt.Errorf("failed to parse bulk string length: %s", err)
	}

	readBytes := make([]byte, count+2)
	if _, err := io.ReadFull(byteStream, readBytes); err != nil {
		return nil, fmt.Errorf("failed to read bulk string contents: %s", err)
	}

	return &Value{
		typ:   BulkString,
		bytes: readBytes[:count],
	}, nil
}

func decodeArray(byteStream *bufio.Reader) (*Value, error) {
	readBytesForCount, err := readUntilCRLF(byteStream)
	if err != nil {
		return nil, fmt.Errorf("failed to read bulk string length: %s", err)
	}
	
	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return nil, fmt.Errorf("failed to parse bulk string length: %s", err)
	}

	array := []Value{}
	for i := 0; i < count; i++ {
		value, err := DecodeRESP(byteStream)
		if err != nil {
			return nil, err
		}
		array = append(array, *value)
	}

	return &Value{
		typ:     Array,
		command: array[0].String(),
		args:    array[1:],
	}, nil
}

func prepareRESPString(s string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", s))
}

func prepareRESPArray(args []Value) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0].String()), args[0].String()))
}
