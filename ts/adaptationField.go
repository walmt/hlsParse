package ts

import (
	"fmt"
	"hlsParse/util"
)

const (
	DiscontinuityIndicatorMark            = 0b10000000
	RandomAccessIndicatorMark             = 0b01000000
	ElementaryStreamPriorityIndicatorMark = 0b00100000
	PcrFlagMark                           = 0b00010000
	OpcrFlagMark                          = 0b00001000
	SplicingPointFlagMark                 = 0b00000100
	TransportPrivateDataFlagMark          = 0b00000010
	AdaptationFieldExtensionFlagMark      = 0b00000001
)

type AdaptationField struct {
	TransportStream *TransportStream

	HaveAdaptationField          bool
	PcrFlag                      uint8
	OpcrFlag                     uint8
	SplicingPointFlag            uint8
	TransportPrivateDataFlag     uint8
	AdaptationFieldExtensionFlag uint8
}

func (a *AdaptationField) Parse(buf []byte, index int) (int, error) {

	adaptationFieldControl := a.TransportStream.TsHeader.AdaptationFieldControl
	if adaptationFieldControl != AdaptationFieldControl10 &&
		adaptationFieldControl != AdaptationFieldControl11 {
		return index, nil
	}
	a.HaveAdaptationField = true

	adaptationFieldLength, err := util.BytesToUint16ByBigEndian(buf[index : index+1])
	if err != nil {
		return 0, fmt.Errorf("util.BytesToUint16ByBigEndian(buf[index : index+1]) failed, err:%v", err)
	}
	fmt.Printf("adaptationFieldLength is %v\n", adaptationFieldLength)

	index += 1

	discontinuityIndicator := (buf[index] & DiscontinuityIndicatorMark) >> 7
	fmt.Printf("discontinuityIndicator is %v\n", discontinuityIndicator)

	randomAccessIndicator := (buf[index] & RandomAccessIndicatorMark) >> 6
	fmt.Printf("randomAccessIndicator is %v\n", randomAccessIndicator)

	elementaryStreamPriorityIndicator := (buf[index] & ElementaryStreamPriorityIndicatorMark) >> 5
	fmt.Printf("elementaryStreamPriorityIndicator is %v\n", elementaryStreamPriorityIndicator)

	pcrFlag := (buf[index] & PcrFlagMark) >> 4
	a.PcrFlag = pcrFlag
	fmt.Printf("pcrFlag is %v\n", pcrFlag)

	opcrFlag := (buf[index] & OpcrFlagMark) >> 3
	a.OpcrFlag = opcrFlag
	fmt.Printf("OpcrFlag us %v\n", opcrFlag)

	splicingPointFlag := (buf[index] & SplicingPointFlagMark) >> 2
	a.SplicingPointFlag = splicingPointFlag
	fmt.Printf("splicingPointFlag is %v\n", splicingPointFlag)

	transportPrivateDataFlag := (buf[index] & TransportPrivateDataFlagMark) >> 1
	a.TransportPrivateDataFlag = transportPrivateDataFlag
	fmt.Printf("transportPrivateDataFlag is %v\n", transportPrivateDataFlag)

	adaptationFieldExtensionFlag := (buf[index] & AdaptationFieldExtensionFlagMark) >> 0
	a.AdaptationFieldExtensionFlag = adaptationFieldControl
	fmt.Printf("adaptationFieldExtensionFlag is %v\n", adaptationFieldExtensionFlag)

	index += 1

	if pcrFlag == 1 {

	}

	return 0, nil
}
