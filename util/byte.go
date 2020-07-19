package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func BytesIsZero(buf []byte) bool {
	for i := 0; i < len(buf); i++ {
		if buf[i] != 0 {
			return false
		}
	}
	return true
}

func BytesIsOxFF(buf []byte) bool {
	for i := 0; i < len(buf); i++ {
		if buf[i] != 0xFF {
			return false
		}
	}
	return true
}

func BytesToUint64ByBigEndian(buf []byte) (uint64, error) {

	if len(buf) > 8 {
		return 0, fmt.Errorf("buf len more than 8, len:%v", len(buf))
	}

	b := make([]byte, 8)
	less := 8 - len(buf)
	for i := 0; i < len(buf); i++ {
		b[i+less] = buf[i]
	}

	bytesBuffer := bytes.NewBuffer(b)

	var x uint64
	if err := binary.Read(bytesBuffer, binary.BigEndian, &x); err != nil {
		return 0, fmt.Errorf("binary.Read failed, err:%v", err)
	}

	return x, nil
}

func BytesToInt32ByBigEndian(buf []byte) (int32, error) {

	if len(buf) > 4 {
		return 0, fmt.Errorf("buf len more than 4, len:%v", len(buf))
	}

	b := make([]byte, 4)
	less := 4 - len(buf)
	for i := 0; i < len(buf); i++ {
		b[i+less] = buf[i]
	}

	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	if err := binary.Read(bytesBuffer, binary.BigEndian, &x); err != nil {
		return 0, fmt.Errorf("binary.Read failed, err:%v", err)
	}

	return x, nil
}

func BytesToUint32ByBigEndian(buf []byte) (uint32, error) {

	if len(buf) > 4 {
		return 0, fmt.Errorf("buf len more than 4, len:%v", len(buf))
	}

	b := make([]byte, 4)
	less := 4 - len(buf)
	for i := 0; i < len(buf); i++ {
		b[i+less] = buf[i]
	}

	bytesBuffer := bytes.NewBuffer(b)

	var x uint32
	if err := binary.Read(bytesBuffer, binary.BigEndian, &x); err != nil {
		return 0, fmt.Errorf("binary.Read failed, err:%v", err)
	}

	return x, nil
}

func BytesToUint16ByBigEndian(buf []byte) (uint16, error) {

	if len(buf) > 2 {
		return 0, fmt.Errorf("buf len more than 2, len:%v", len(buf))
	}

	b := make([]byte, 2)
	less := 2 - len(buf)
	for i := 0; i < len(buf); i++ {
		b[i+less] = buf[i]
	}

	bytesBuffer := bytes.NewBuffer(b)

	var x uint16
	if err := binary.Read(bytesBuffer, binary.BigEndian, &x); err != nil {
		return 0, fmt.Errorf("binary.Read failed, err:%v", err)
	}

	return x, nil
}

func BytesToUint8ByBigEndian(buf byte) uint8 {

	b := make([]byte, 1)
	b[0] = buf

	bytesBuffer := bytes.NewBuffer(b)

	var x uint8
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)

	return x
}

func ByteToFloat64(buf []byte) (float64, error) {

	if len(buf) > 8 {
		return 0, fmt.Errorf("buf len more than 8, len:%v", len(buf))
	}

	b := make([]byte, 8)
	less := 8 - len(buf)
	for i := 0; i < len(buf); i++ {
		b[i+less] = buf[i]
	}
	//
	//bits := binary.BigEndian.Uint64(buf)
	//
	//return math.Float64frombits(bits), nil

	bytesBuffer := bytes.NewBuffer(b)

	var x float64
	if err := binary.Read(bytesBuffer, binary.BigEndian, &x); err != nil {
		return 0, fmt.Errorf("binary.Read failed, err:%v", err)
	}

	return x, nil
}
