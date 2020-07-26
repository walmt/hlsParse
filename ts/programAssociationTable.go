package ts

import (
	"fmt"
	"hlsParse/util"
)

const (
	SectionSyntaxIndicatorMark = 0b10000000
	ZeroMark                   = 0b01000000
	Reserved0Mark              = 0b00110000
	SectionLengthMark          = 0b00001111
	Reserved1Mark              = 0b11000000
	VersionNumberMark          = 0b00111110
	CurrentNextIndicatorMark   = 0b00000001
	Reserved2Mark              = 0b11100000
	NetworkIdMark              = 0b00011111
	ProgramMapPidMark          = 0b00011111
)

const (
	CurrentNextIndicator0 = 0
	CurrentNextIndicator1 = 1

	ProgramNumber0 = "0x0000"
	ProgramNumber1 = "0x0001"
)

var CurrentNextIndicatorMap = map[uint8]string{
	CurrentNextIndicator0: "0, 表示下一个表有效",
	CurrentNextIndicator1: "1, 表示传送的PAT当前可以使用",
}

var ProgramNumberMap = map[string]string{
	ProgramNumber0: "0x0000, 后面的PID是网络PID",
	ProgramNumber1: "0x0001, 这个为PMT",
}

type ProgramAssociationTable struct {
	Payload *Payload
}

func (p *ProgramAssociationTable) Parse(buf []byte, index int) (int, error) {

	// 去掉调整字节
	if p.Payload.TransportStream.TsHeader.PayloadUnitStartIndicator == 1 {
		index += 1
	}

	tableId := buf[index]
	if tableId != 0 {
		return 0, fmt.Errorf("tableId != 0")
	}
	fmt.Printf("tableId is 0\n")

	index += 1

	sectionSyntaxIndicator := (buf[index] & SectionSyntaxIndicatorMark) >> 7
	if sectionSyntaxIndicator != 1 {
		return 0, fmt.Errorf("sectionSyntaxIndicator is not 1")
	}
	fmt.Printf("sectionSyntaxIndicator is %v\n", sectionSyntaxIndicator)

	isZero := (buf[index] & ZeroMark) >> 6
	if isZero != 0 {
		return 0, fmt.Errorf("isZero != 0")
	}
	fmt.Printf("isZero is %v\n", isZero)

	reserved0 := (buf[index] & Reserved0Mark) >> 4
	if reserved0 != 0b11 {
		return 0, fmt.Errorf("reserved0 != 0b11")
	}
	fmt.Printf("reserved0 is 0b%b\n", reserved0)

	sectionLengthBuf := buf[index : index+2]
	sectionLengthBuf[0] = sectionLengthBuf[0] & SectionLengthMark
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

	reserved1 := (buf[index] & Reserved1Mark) >> 6
	if reserved1 != 0b11 {
		return 0, fmt.Errorf("reserved1 != 0b11")
	}
	fmt.Printf("reserved1 is 0b%b\n", reserved1)

	versionNumber := (buf[index] & VersionNumberMark) >> 1
	fmt.Printf("versionNumber is %v\n", versionNumber)

	currentNextIndicator := buf[index] & CurrentNextIndicatorMark
	fmt.Printf("currentNextIndicator is %v\n", CurrentNextIndicatorMap[currentNextIndicator])

	index += 1

	sectionNumber := buf[index]
	fmt.Printf("sectionNumber is %v\n", sectionNumber)
	index += 1

	lastSectionNumber := buf[index]
	fmt.Printf("lastSectionNumber is %v\n", lastSectionNumber)

	index += 1

	for i := 0; i < int(sectionLength)-12; i += 4 {
		programNumber := buf[index : index+2]
		programNumberStr, ok := ProgramNumberMap[fmt.Sprintf("0x%x", programNumber)]
		if !ok {
			programNumberStr = fmt.Sprintf("0x%x, 其他值，由用户自定义", programNumberStr)
		}
		fmt.Printf("programNumber is %v\n", programNumberStr)

		index += 2

		reserved2 := (buf[index] & Reserved2Mark) >> 5
		if reserved2 != 0b111 {
			return 0, fmt.Errorf("reserved2 != 0b111")
		}
		fmt.Printf("reserved2 is 0b111\n")

		if programNumber[1] == 0x00 {

			networkIdBuf := buf[index : index+2]
			networkIdBuf[0] = networkIdBuf[0] & NetworkIdMark
			networkId, err := util.BytesToUint16ByBigEndian(networkIdBuf)
			if err != nil {
				return 0, fmt.Errorf("util.BytesToUint16ByBigEndian(networkIdBuf) is fail err:%v", err)
			}
			fmt.Printf("networkId is %v\n", networkId)

		} else {

			programMapPidBuf := buf[index : index+2]
			programMapPidBuf[0] = programMapPidBuf[0] & NetworkIdMark
			programMapPid, err := util.BytesToUint16ByBigEndian(programMapPidBuf)
			if err != nil {
				return 0, fmt.Errorf("util.BytesToUint16ByBigEndian(programMapPidBuf) is fail err:%v", err)
			}
			fmt.Printf("programMapPid is %v\n", programMapPid)

		}

		index += 2
	}

	crc32 := buf[index : index+4]
	fmt.Printf("crc32 is 0x%x\n", crc32)

	index += 4
	
	fmt.Println()

	return index, nil
}
