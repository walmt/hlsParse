package ts

import "fmt"

const (
	TransportErrorIndicatorMark    = 0b10000000
	PayloadUnitStartIndicatorMark  = 0b01000000
	TransportPriorityMark          = 0b00100000
	PidPreMark                     = 0b00011111
	TransportScramblingControlMark = 0b11000000
	AdaptationFieldControlMark     = 0b00110000
	ContinuityCounterMark          = 0b00001111
)

const (
	PayloadUnitStartIndicator0 = 0
	PayloadUnitStartIndicator1 = 1

	AdaptationFieldControl00 = 0b00000000 // 保留值
	AdaptationFieldControl01 = 0b00000001 // 负载中只有有效载荷
	AdaptationFieldControl10 = 0b00000010 // 负载中只有自适应字段
	AdaptationFieldControl11 = 0b00000011 // 先有自适应字段，再有有效载荷

	AdaptationFieldFlagHavePCR = 0x50
	AdaptationFieldFlagNoPCR   = 0x40
)

var PayloadUnitStartIndicatorMap = map[uint8]string{
	PayloadUnitStartIndicator0: "0 - 如果有PES包，表示TS包的开始不是PES包；如果带有PSI，在有效净荷中没有指针pointer_field",
	PayloadUnitStartIndicator1: "1 - 如果有PES包，表示TS包的有效净荷以PES包的第一个字节开始；如果带有PSI，即第一个字节带有指针pointer_field。",
}

var AdaptationFieldControlMap = map[uint8]string{
	AdaptationFieldControl00: "0b00（保留值）",
	AdaptationFieldControl01: "0b01（负载中只有有效载荷）",
	AdaptationFieldControl10: "0b10（负载中只有自适应字段）",
	AdaptationFieldControl11: "0b11（先有自适应字段，再有有效载荷）",
}

var AdaptationFieldFlagMap = map[byte]string{
	AdaptationFieldFlagNoPCR:   "0x40（没有PCR）",
	AdaptationFieldFlagHavePCR: "0x50（有PCR）",
}

type TsHeader struct {
	TransportStream           *TransportStream
	Pid                       uint16
	AdaptationFieldControl    uint8
	PayloadUnitStartIndicator byte
}

func (t *TsHeader) Parse(buf []byte, index int) (int, error) {
	if len(buf[index:]) < 4 { // header长度固定为4
		return 0, fmt.Errorf("len(buf) != 4")
	}

	syncByte := buf[index]
	if syncByte != 0x47 {
		return 0, fmt.Errorf("syncByte != 0x47")
	}
	fmt.Printf("syncByte is 0x47\n")

	index += 1

	transportErrorIndicator := (buf[index] & TransportErrorIndicatorMark) >> 7
	if transportErrorIndicator != 0 {
		return 0, fmt.Errorf("transportErrorIndicator != 0")
	}
	fmt.Printf("TransportErrorIndicator is 0\n")

	payloadUnitStartIndicator := (buf[index] & PayloadUnitStartIndicatorMark) >> 6
	payloadUnitStartIndicatorString, ok := PayloadUnitStartIndicatorMap[payloadUnitStartIndicator]
	if !ok {
		return 0, fmt.Errorf("PayloadUnitStartIndicatorMap[payloadUnitStartIndicator] failed, "+
			"payloadUnitStartIndicator:%v", payloadUnitStartIndicator)
	}
	t.PayloadUnitStartIndicator = payloadUnitStartIndicator
	fmt.Printf("PayloadUnitStartIndicator is %v\n", payloadUnitStartIndicatorString)

	transportPriority := (buf[index] & TransportPriorityMark) >> 6
	fmt.Printf("TransportPriority is %v\n", transportPriority)

	pid := uint16(buf[index]&PidPreMark)<<8 + uint16(buf[index+1])
	fmt.Printf("Pid is 0x%x\n", pid)

	t.Pid = pid
	index += 2

	transportScramblingControl := (buf[index] & TransportScramblingControlMark) >> 6
	if transportScramblingControl != 0 {
		return 0, fmt.Errorf("transportScramblingControl != 0")
	}
	fmt.Printf("TransportScramblingControl is %v\n", transportScramblingControl)

	adaptationFieldControl := (buf[index] & AdaptationFieldControlMark) >> 4
	adaptationFieldControlString, ok := AdaptationFieldControlMap[adaptationFieldControl]
	if !ok {
		return 0, fmt.Errorf("AdaptationFieldControlMap[adaptationFieldControl] failed, "+
			"adaptationFieldControl:%b", adaptationFieldControl)
	}
	t.AdaptationFieldControl = adaptationFieldControl
	fmt.Printf("adaptationFieldControl is %v\n", adaptationFieldControlString)

	continuityCounter := (buf[index] & ContinuityCounterMark) >> 0
	fmt.Printf("ContinuityCounter is %v\n", continuityCounter)

	index += 1

	fmt.Println()

	return index, nil
}
