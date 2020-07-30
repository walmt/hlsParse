package ts

import (
	"fmt"
	"hlsParse/util"
)

type ProgramMapTable struct {
	Ts               *Ts
	PidStreamTypeMap map[uint16]uint8
}

const (
	PmtSectionSyntaxIndicatorMark = 0b10000000
	PmtZeroMark                   = 0b01000000
	PmtReserved0Mark              = 0b00110000
	PmtSectionLengthMark          = 0b00001111
	PmtReserved1Mark              = 0b11000000
	PmtVersionNumberMark          = 0b00111110
	PmtCurrentNextIndicatorMark   = 0b00000001
	PmtReserved2Mark              = 0b11100000
	PmtPcrPidMark                 = 0b00011111
	PmtReserved3Mark              = 0b11110000
	PmtProgramInfoLengthMark      = 0b00001111
	PmtReserved4Mark              = 0b11100000
	PmtElementaryPidMark          = 0b00011111
	PmtReserved5Mark              = 0b11110000
	PmtEsInfoLengthMark           = 0b00001111
)

const (
	PmtCurrentNextIndicator0 = 0
	PmtCurrentNextIndicator1 = 1
)

var PmtCurrentNextIndicatorMap = map[uint8]string{
	PmtCurrentNextIndicator0: "0, 表示当前传送的program_map_section不可用，下一个TS的program_map_section有效。",
	PmtCurrentNextIndicator1: "1, 表示当前传送的program_map_section可用",
}

func BuildProgramMapTable(t *Ts) *ProgramMapTable {

	p := new(ProgramMapTable)
	p.PidStreamTypeMap = make(map[uint16]uint8)
	p.Ts = t

	return p
}

func (p *ProgramMapTable) Parse(buf []byte, index int) (int, error) {

	// 去掉调整字节
	if p.Ts.CurrentTransportStream.TsHeader.PayloadUnitStartIndicator == 1 {
		index += 1
	}

	tableId := buf[index]
	if tableId != 0x02 {
		return 0, fmt.Errorf("tableId != 0x02")
	}
	fmt.Printf("tableId is 0x02\n")

	index += 1

	sectionSyntaxIndicator := (buf[index] & PmtSectionSyntaxIndicatorMark) >> 7
	if sectionSyntaxIndicator != 1 {
		return 0, fmt.Errorf("sectionSyntaxIndicator is not 1")
	}
	fmt.Printf("sectionSyntaxIndicator is %v\n", sectionSyntaxIndicator)

	isZero := (buf[index] & PmtZeroMark) >> 6
	if isZero != 0 {
		return 0, fmt.Errorf("isZero != 0")
	}
	fmt.Printf("isZero is 0\n")

	reserved0 := (buf[index] & PmtReserved0Mark) >> 4
	if reserved0 != 0b11 {
		return 0, fmt.Errorf("reserved0 != 0b11")
	}
	fmt.Printf("reserved0 is 0b%b\n", reserved0)

	sectionLengthBuf := buf[index : index+2]
	sectionLengthBuf[0] = sectionLengthBuf[0] & PmtSectionLengthMark
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

	reserved1 := (buf[index] & PmtReserved1Mark) >> 6
	if reserved1 != 0b11 {
		return 0, fmt.Errorf("reserved1 != 0b11")
	}
	fmt.Printf("reserved1 is 0b%b\n", reserved1)

	versionNumber := (buf[index] & PmtVersionNumberMark) >> 1
	fmt.Printf("versionNumber is %v\n", versionNumber)

	currentNextIndicator := buf[index] & PmtCurrentNextIndicatorMark
	fmt.Printf("currentNextIndicator is %v\n", PmtCurrentNextIndicatorMap[currentNextIndicator])

	index += 1

	sectionNumber := buf[index]
	fmt.Printf("sectionNumber is %v\n", sectionNumber)
	index += 1

	lastSectionNumber := buf[index]
	fmt.Printf("lastSectionNumber is %v\n", lastSectionNumber)

	index += 1

	reserved2 := (buf[index] & PmtReserved2Mark) >> 5
	if reserved2 != 0b111 {
		return 0, fmt.Errorf("reserved2 != 0b111")
	}
	fmt.Printf("reserved2 is 0b111\n")

	pcrPidBuf := buf[index : index+2]
	pcrPidBuf[0] = pcrPidBuf[0] & PmtPcrPidMark
	pcrPid, err := util.BytesToUint16ByBigEndian(pcrPidBuf)
	if err != nil {
		return 0, fmt.Errorf("util.BytesToUint16ByBigEndian(pcrPidBuf) failed, err:%v", err)
	}
	fmt.Printf("pcrPid is 0x%x\n", pcrPid)

	index += 2

	reserved3 := (buf[index] & PmtReserved3Mark) >> 4
	if reserved3 != 0b1111 {
		return 0, fmt.Errorf("reserved3 != 0b1111")
	}
	fmt.Printf("reserved3 is 0b1111\n")

	programInfoLengthBuf := buf[index : index+2]
	programInfoLengthBuf[0] = programInfoLengthBuf[0] & PmtProgramInfoLengthMark
	programInfoLength, err := util.BytesToUint16ByBigEndian(programInfoLengthBuf)
	if err != nil {
		return 0, fmt.Errorf("util.BytesToUint16ByBigEndian(programInfoLengthBuf) failed, err:%v", err)
	}
	fmt.Printf("programInfoLength is %v\n", programInfoLength)

	index += 2

	index += int(programInfoLength)

	descriptorIndexEnd := int(sectionLength) - 13 + index

	for index < descriptorIndexEnd {
		streamType := buf[index]
		fmt.Printf("streamType is 0x%02x(0x1b - h264、0x0f - AAC、0x03 - mp3)\n", streamType)

		index += 1

		reserved4 := (buf[index] & PmtReserved4Mark) >> 5
		if reserved4 != 0b111 {
			return 0, fmt.Errorf("reserved4 != 0b111")
		}
		fmt.Printf("reserved4 is 0b111\n")

		elementaryPidBuf := buf[index : index+2]
		elementaryPidBuf[0] = elementaryPidBuf[0] & PmtElementaryPidMark
		elementaryPid, err := util.BytesToUint16ByBigEndian(elementaryPidBuf)
		if err != nil {
			return 0, fmt.Errorf("util.BytesToUint16ByBigEndian(elementaryPidBuf) failed, err:%v", err)
		}
		p.PidStreamTypeMap[elementaryPid] = streamType
		fmt.Printf("elementaryPid is 0x%x\n", elementaryPid)

		index += 2

		reserved5 := (buf[index] & PmtReserved5Mark) >> 4
		if reserved5 != 0b1111 {
			return 0, fmt.Errorf("reserved5 != 0b1111")
		}
		fmt.Printf("reserved5 is 0b1111\n")

		esInfoLengthBuf := buf[index : index+2]
		esInfoLengthBuf[0] = esInfoLengthBuf[0] & PmtEsInfoLengthMark
		esInfoLength, err := util.BytesToUint16ByBigEndian(esInfoLengthBuf)
		if err != nil {
			return 0, fmt.Errorf("util.BytesToUint16ByBigEndian(esInfoLengthBuf) failed, err:%v", err)
		}
		fmt.Printf("esInfoLength is %v\n", esInfoLength)
		index += int(esInfoLength)

		index += 2
	}

	crc32 := buf[index : index+4]
	fmt.Printf("crc32 is 0x%x\n", crc32)

	index += 4

	fmt.Println()

	return index, nil
}
