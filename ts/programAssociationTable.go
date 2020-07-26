package ts

import (
	"fmt"
	"hlsParse/util"
)

const (
	PatSectionSyntaxIndicatorMark = 0b10000000
	PatZeroMark                   = 0b01000000
	PatReserved0Mark              = 0b00110000
	PatSectionLengthMark          = 0b00001111
	PatReserved1Mark              = 0b11000000
	PatVersionNumberMark          = 0b00111110
	PatCurrentNextIndicatorMark   = 0b00000001
	PatReserved2Mark              = 0b11100000
	PatNetworkIdMark              = 0b00011111
	PatProgramMapPidMark          = 0b00011111
)

const (
	PatCurrentNextIndicator0 = 0
	PatCurrentNextIndicator1 = 1

	PatProgramNumber0 = "0x0000"
	PatProgramNumber1 = "0x0001"
)

var PatCurrentNextIndicatorMap = map[uint8]string{
	PatCurrentNextIndicator0: "0, 表示下一个表有效",
	PatCurrentNextIndicator1: "1, 表示传送的PAT当前可以使用",
}

var PatProgramNumberMap = map[string]string{
	PatProgramNumber0: "0x0000, 后面的PID是网络PID",
	PatProgramNumber1: "0x0001, 这个为PMT",
}

type ProgramAssociationTable struct {
	Ts            *Ts
	ProgramMapPid uint16
}

func (p *ProgramAssociationTable) Parse(buf []byte, index int) (int, error) {

	// 去掉调整字节
	if p.Ts.CurrentTransportStream.TsHeader.PayloadUnitStartIndicator == 1 {
		index += 1
	}

	tableId := buf[index]
	if tableId != 0 {
		return 0, fmt.Errorf("tableId != 0")
	}
	fmt.Printf("tableId is 0\n")

	index += 1

	sectionSyntaxIndicator := (buf[index] & PatSectionSyntaxIndicatorMark) >> 7
	if sectionSyntaxIndicator != 1 {
		return 0, fmt.Errorf("sectionSyntaxIndicator is not 1")
	}
	fmt.Printf("sectionSyntaxIndicator is %v\n", sectionSyntaxIndicator)

	isZero := (buf[index] & PatZeroMark) >> 6
	if isZero != 0 {
		return 0, fmt.Errorf("isZero != 0")
	}
	fmt.Printf("isZero is %v\n", isZero)

	reserved0 := (buf[index] & PatReserved0Mark) >> 4
	if reserved0 != 0b11 {
		return 0, fmt.Errorf("reserved0 != 0b11")
	}
	fmt.Printf("reserved0 is 0b%b\n", reserved0)

	sectionLengthBuf := buf[index : index+2]
	sectionLengthBuf[0] = sectionLengthBuf[0] & PatSectionLengthMark
	sectionLength, err := util.BytesToUint16ByBigEndian(sectionLengthBuf)
	if err != nil {
		return 0, fmt.Errorf("util.BytesToUint16ByBigEndian(sectionLengthBuf) failed, err:%v", err)
	}
	fmt.Printf("sectionLength is %v\n", sectionLength)

	index += 2

	transportStreamId, err := util.BytesToUint16ByBigEndian(buf[index : index+2])
	if err != nil {
		return 0, fmt.Errorf("util.BytesToUint16ByBigEndian(buf[index : index+2]) failed, err:%v", err)
	}
	fmt.Printf("transportStreamId is %v\n", transportStreamId)

	index += 2

	reserved1 := (buf[index] & PatReserved1Mark) >> 6
	if reserved1 != 0b11 {
		return 0, fmt.Errorf("reserved1 != 0b11")
	}
	fmt.Printf("reserved1 is 0b%b\n", reserved1)

	versionNumber := (buf[index] & PatVersionNumberMark) >> 1
	fmt.Printf("versionNumber is %v\n", versionNumber)

	currentNextIndicator := buf[index] & PatCurrentNextIndicatorMark
	fmt.Printf("currentNextIndicator is %v\n", PatCurrentNextIndicatorMap[currentNextIndicator])

	index += 1

	sectionNumber := buf[index]
	fmt.Printf("sectionNumber is %v\n", sectionNumber)
	index += 1

	lastSectionNumber := buf[index]
	fmt.Printf("lastSectionNumber is %v\n", lastSectionNumber)

	index += 1

	for i := 0; i < int(sectionLength)-12; i += 4 {
		programNumber := buf[index : index+2]
		programNumberStr, ok := PatProgramNumberMap[fmt.Sprintf("0x%x", programNumber)]
		if !ok {
			programNumberStr = fmt.Sprintf("0x%x, 其他值，由用户自定义", programNumberStr)
		}
		fmt.Printf("programNumber is %v\n", programNumberStr)

		index += 2

		reserved2 := (buf[index] & PatReserved2Mark) >> 5
		if reserved2 != 0b111 {
			return 0, fmt.Errorf("reserved2 != 0b111")
		}
		fmt.Printf("reserved2 is 0b111\n")

		if programNumber[1] == 0x00 {

			networkIdBuf := buf[index : index+2]
			networkIdBuf[0] = networkIdBuf[0] & PatNetworkIdMark
			networkId, err := util.BytesToUint16ByBigEndian(networkIdBuf)
			if err != nil {
				return 0, fmt.Errorf("util.BytesToUint16ByBigEndian(networkIdBuf) is fail err:%v", err)
			}
			fmt.Printf("networkId is %v\n", networkId)

		} else {

			programMapPidBuf := buf[index : index+2]
			programMapPidBuf[0] = programMapPidBuf[0] & PatProgramMapPidMark
			programMapPid, err := util.BytesToUint16ByBigEndian(programMapPidBuf)
			if err != nil {
				return 0, fmt.Errorf("util.BytesToUint16ByBigEndian(programMapPidBuf) is fail err:%v", err)
			}
			p.ProgramMapPid = programMapPid
			fmt.Printf("programMapPid is 0x%x\n", programMapPid)
		}

		index += 2
	}

	crc32 := buf[index : index+4]
	fmt.Printf("crc32 is 0x%x\n", crc32)

	index += 4

	fmt.Println()

	return index, nil
}
