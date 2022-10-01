package kitus

import (
	"encoding/binary"
	"io"
)

func readHeader(r io.Reader) (int32, error) {
	return readInt32(r)
}

func readInt16(r io.Reader) (int16, error) {
	b := make([]byte, 2)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return 0, err
	}
	return int16(binary.BigEndian.Uint16(b)), nil
}

func readInt32(r io.Reader) (int32, error) {
	b := make([]byte, 4)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(b)), nil
}

func readString(r io.Reader) (string, error) {
	n, err := readInt16(r)
	if err != nil {
		return "", err
	}
	b := make([]byte, n)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
