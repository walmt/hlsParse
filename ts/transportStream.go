package ts

import (
	"fmt"
)

type TransportStream struct {
	TsHeader        *TsHeader
	AdaptationField *AdaptationField
	Payload         *Payload
}

type TsHeader struct {
	TransportStream           *TransportStream
	PayloadUnitStartIndicator uint8
	AdaptationFieldControl    uint8
}

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

type AdaptationField struct {
	TransportStream *TransportStream

	HaveAdaptationField bool
}

type Payload struct {
	TransportStream *TransportStream
}

func (t *TransportStream) Parse(buf []byte) error {

	if len(buf) != 188 {
		return fmt.Errorf("len(buf) != 188")
	}

	var err error
	index := 0
	tsHeader := t.getTsHeader()
	index, err = tsHeader.Parse(buf, index)
	if err != nil {
		return fmt.Errorf("tsHeader.Parse failed, err:%v", err)
	}
	fmt.Println()

	adaptationField := t.getAdaptationField()
	index, err = adaptationField.Parse(buf, index)
	if err != nil {
		return fmt.Errorf("adaptationField.Parse failed, err:%v", err)
	}
	fmt.Println()

	return nil
}

func (t *TransportStream) getTsHeader() *TsHeader {

	if t.TsHeader != nil {
		return t.TsHeader
	}

	t.TsHeader = new(TsHeader)
	t.TsHeader.TransportStream = t

	return t.TsHeader
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

	return index, nil
}

func (t *TransportStream) getAdaptationField() *AdaptationField {
	if t.AdaptationField != nil {
		return t.AdaptationField
	}

	t.AdaptationField = new(AdaptationField)
	t.AdaptationField.TransportStream = t

	return t.AdaptationField
}

func (a *AdaptationField) Parse(buf []byte, index int) (int, error) {

	adaptationFieldControl := a.TransportStream.TsHeader.AdaptationFieldControl
	if adaptationFieldControl != AdaptationFieldControl10 &&
		adaptationFieldControl != AdaptationFieldControl11 {
		return index, nil
	}
	a.HaveAdaptationField = true

	adaptationFieldLength := buf[index : index+8]
	if a.TransportStream.TsHeader.AdaptationFieldControl == AdaptationFieldControl11 {
		// 既有自适应区又有有效荷载时，adaptationFieldLength介于0~182

	}


	index += 8

	return 0, nil
}

//func (a *AdaptationField) Parse(buf []byte, index int) (int, error) {
//
//	adaptationFieldControl := a.TransportStream.TsHeader.AdaptationFieldControl
//	if adaptationFieldControl != AdaptationFieldControl10 &&
//		adaptationFieldControl != AdaptationFieldControl11 {
//		return index, nil
//	}
//
//	a.HaveAdaptationField = true
//	adaptationFieldLength := buf[index]
//	if len(buf[index+1:]) < int(adaptationFieldLength) {
//		return 0, fmt.Errorf("len(buf[index+1:]) < int(adaptationFieldLength)")
//	}
//	fmt.Printf("AdaptationField Length is %v\n", adaptationFieldLength)
//
//	index += 1
//
//	flag := buf[index]
//	flagString, ok := AdaptationFieldFlagMap[flag]
//	if !ok {
//		return 0, fmt.Errorf("AdaptationFieldFlagMap[flag] failed, flag:%x", flag)
//	}
//	fmt.Printf("AdaptationField Flag is %v\n", flagString)
//
//	index += 1
//
//	pcr := buf[index : index+5]
//	if flag == AdaptationFieldFlagNoPCR && !util.BytesIsZero(pcr) {
//		return 0, fmt.Errorf("have no pcr but pcr is not zero")
//	}
//	fmt.Printf("AdaptationField PCR is %v\n", pcr)
//
//	index += 5
//
//	stuffingBytes := buf[index : index+int(adaptationFieldLength)-6]
//	if !util.BytesIsOxFF(stuffingBytes) {
//		return 0, fmt.Errorf("AdaptationField StuffingBytes is not 0xFF, StuffingBytes:%v", stuffingBytes)
//	}
//	fmt.Printf("AdaptationField StuffingBytes length:%v\n", len(stuffingBytes))
//
//	index += int(adaptationFieldLength) - 6
//
//	return index, nil
//}

func (t *TransportStream) getPayload() *Payload {
	if t.Payload != nil {
		return t.Payload
	}

	t.Payload = new(Payload)
	t.Payload.TransportStream = t

	return t.Payload
}
