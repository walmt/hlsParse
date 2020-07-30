package ts

import (
	"fmt"
	"hlsParse/util"
	"os"
	"strconv"
)

type ElementaryStreams struct {
	Pid                         uint16
	StreamId                    uint8
	PacketizedElementaryStreams *PacketizedElementaryStreams
	DataFile                    *os.File
}

func buildElementaryStream(p *PacketizedElementaryStreams) *ElementaryStreams {

	e := new(ElementaryStreams)
	e.PacketizedElementaryStreams = p
	e.Pid = p.Ts.CurrentTransportStream.TsHeader.Pid

	return e
}

func (e *ElementaryStreams) ParseHeader(buf []byte, index int) (int, error) {
	packetStartCodePrefix, err := util.BytesToUint32ByBigEndian(buf[index : index+3])
	if err != nil {
		return 0, fmt.Errorf("util.BytesToUint32ByBigEndian(buf[index : index+3]) failed, err:%v", err)
	}
	if packetStartCodePrefix != 0x000001 {
		return 0, fmt.Errorf("packetStartCodePrefix != 0x000001, packetStartCodePrefix:0x%06x", packetStartCodePrefix)
	}
	fmt.Printf("packetStartCodePrefix is 0x%x\n", packetStartCodePrefix)

	index += 3

	streamId := buf[index]
	e.StreamId = streamId
	fmt.Printf("streamId is 0x%x.(音频取值（0xc0-0xdf），通常为0xc0;视频取值（0xe0-0xef），通常为0xe0)\n", streamId)

	index += 1

	pesPacketLength, err := util.BytesToUint16ByBigEndian(buf[index : index+2])
	if err != nil {
		return 0, fmt.Errorf("util.BytesToUint16ByBigEndian(buf[index : index+2]) failed, err:%v", err)
	}
	fmt.Printf("pesPacketLength is %v\n", pesPacketLength)
	
	index += 2

	flag0 := buf[index]
	fmt.Printf("flag0 is 0x%x.(通常取值0x80，表示数据不加密、无优先级、备份的数据)\n", flag0)

	index += 1

	flag1 := buf[index]
	if flag1 != 0x80 && flag1 != 0xc0 {
		return 0, fmt.Errorf("not support parse flag1(0x%x)", flag1)
	}
	fmt.Printf("flag1 is 0x%x.(取值0x80表示只含有pts，取值0xc0表示含有pts和dts)\n", flag1)

	index += 1

	pesDataLength := buf[index]
	if (flag1 == 0x80 && pesDataLength != 5) || (flag1 == 0xc0 && pesDataLength != 10) {
		return 0, fmt.Errorf("(flag1 == 0x80 && pesDataLength != 5) || (flag1 == 0xc0 && pesDataLength != 10), flag1:0x%x, pesDataLength:%v", flag1, pesDataLength)
	}
	fmt.Printf("pesDataLength is %v\n", pesDataLength)

	index += 1

	ptsConstValue0 := (buf[index] & PtsConstValue0Mark) >> 4
	if ptsConstValue0 != 0b0011 && ptsConstValue0 != 0b0010 {
		return 0, fmt.Errorf("ptsConstValue0 != 0b0011 && ptsConstValue0 != 0b0010")
	}
	fmt.Printf("ptsConstValue0 is 0b%08b\n", ptsConstValue0)

	pts0 := buf[index] & PtsPtsAMark

	ptsConstValue1 := buf[index] & PtsConstValue1Mark
	if ptsConstValue1 != 1 {
		return 0, fmt.Errorf("ptsConstValue1 != 1")
	}
	fmt.Printf("ptsConstValue1 is 0b%08b\n", ptsConstValue1)

	index += 1

	pts1 := buf[index]

	index += 1

	pts2 := buf[index] & PtsPtsBMark

	ptsConstValue2 := buf[index] & PtsConstValue2Mark
	if ptsConstValue2 != 0b1 {
		return 0, fmt.Errorf("ptsConstValue2 != 1")
	}
	fmt.Printf("ptsConstValue2 is 0b%08b\n", ptsConstValue2)

	index += 1

	pts3 := buf[index]

	index += 1

	pts4 := buf[index] & PtsPtsCMark
	ptsConstValue3 := buf[index] & PtsConstValue3Mark
	if ptsConstValue3 != 0b1 {
		return 0, fmt.Errorf("ptsConstValue3 != 1")
	}
	fmt.Printf("ptsConstValue3 is 0b%08b\n", ptsConstValue3)

	ptsBuf := make([]byte, 5)
	ptsBuf[4] = pts4>>1 + pts3<<7
	ptsBuf[3] = pts3>>1 + pts2<<6
	ptsBuf[2] = pts2>>2 + pts1<<6
	ptsBuf[1] = pts1>>2 + pts0<<5
	ptsBuf[0] = pts0 >> 3

	pts, err := util.BytesToUint64ByBigEndian(ptsBuf)
	if err != nil {
		fmt.Printf("util.BytesToUint64ByBigEndian(ptsBuf) failed, err:%v", err)
	}
	fmt.Printf("pts is %v\n", pts)

	index += 1

	if flag1 == 0xc0 {

		dtsConstValue0 := (buf[index] & DtsConstValue0Mark) >> 4
		if dtsConstValue0 != 0b0001 {
			return 0, fmt.Errorf("dtsConstValue0 != 0b0001")
		}
		fmt.Printf("dtsConstValue0 is 0b%08b\n", dtsConstValue0)

		dts0 := buf[index] & DtsDtsAMark

		dtsConstValue1 := buf[index] & DtsConstValue1Mark
		if dtsConstValue1 != 1 {
			return 0, fmt.Errorf("dtsConstValue1 != 1")
		}
		fmt.Printf("dtsConstValue1 is 0b%08b\n", dtsConstValue1)

		index += 1

		dts1 := buf[index]

		index += 1

		dts2 := buf[index] & DtsDtsBMark

		dtsConstValue2 := buf[index] & DtsConstValue2Mark
		if dtsConstValue2 != 0b1 {
			return 0, fmt.Errorf("dtsConstValue2 != 1")
		}
		fmt.Printf("dtsConstValue2 is 0b%08b\n", dtsConstValue2)

		index += 1

		dts3 := buf[index]

		index += 1

		dts4 := buf[index] & DtsDtsCMark
		dtsConstValue3 := buf[index] & DtsConstValue3Mark
		if dtsConstValue3 != 0b1 {
			return 0, fmt.Errorf("dtsConstValue3 != 1")
		}
		fmt.Printf("dtsConstValue3 is 0b%08b\n", dtsConstValue3)

		dtsBuf := make([]byte, 5)
		dtsBuf[4] = dts4>>1 + dts3<<7
		dtsBuf[3] = dts3>>1 + dts2<<6
		dtsBuf[2] = dts2>>2 + dts1<<6
		dtsBuf[1] = dts1>>2 + dts0<<5
		dtsBuf[0] = dts0 >> 3

		dts, err := util.BytesToUint64ByBigEndian(dtsBuf)
		if err != nil {
			fmt.Printf("util.BytesToUint64ByBigEndian(dtsBuf) failed, err:%v", err)
		}
		fmt.Printf("dts is %v\n", dts)

		index += 1
	}

	fmt.Println()

	return index, nil
}
func (e *ElementaryStreams) ParseBody(buf []byte, index int) (int, error) {

	file, err := e.getFile()
	if err != nil {
		return 0, fmt.Errorf("e.getFile() failed, err:%v", err)
	}
	_, err = file.Write(buf[index:])
	if err != nil {
		return 0, fmt.Errorf("e.DataFile.Write(buf[index:]) failed, err:%v", err)
	}
	fmt.Printf("data is %x\n", buf[index:])

	return len(buf), nil
}

func (e *ElementaryStreams) getFile() (*os.File, error) {
	if e.DataFile == nil {
		fileName := strconv.FormatInt(int64(e.Pid), 10)
		if e.StreamId >= 0xc0 && e.StreamId <= 0xdf {
			fileName += ".aac"
		}
		if e.StreamId >= 0xe0 && e.StreamId <= 0xef {
			fileName += ".h264"
		}
		dataFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0)
		if err != nil {
			return nil, fmt.Errorf("os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0) failed, fileName:%v, err:%v\n", fileName, err)
		}
		e.DataFile = dataFile
	}

	return e.DataFile, nil
}
