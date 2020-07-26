package ts

import (
	"fmt"
	"hlsParse/util"
)

const (
	DiscontinuityIndicatorMark                 = 0b10000000
	RandomAccessIndicatorMark                  = 0b01000000
	ElementaryStreamPriorityIndicatorMark      = 0b00100000
	PcrFlagMark                                = 0b00010000
	OpcrFlagMark                               = 0b00001000
	SplicingPointFlagMark                      = 0b00000100
	TransportPrivateDataFlagMark               = 0b00000010
	AdaptationFieldExtensionFlagMark           = 0b00000001
	Const1Value0Mark                           = 0b01111110
	ProgramClockReferenceExtensionMark         = 0b00000001
	Const1Value2Mark                           = 0b01111110
	OriginalProgramClockReferenceExtensionMark = 0b00000001
	LtwFlagMark                                = 0b10000000
	PiecewiseRateFlagMark                      = 0b01000000
	SeamlessSpliceFlagMark                     = 0b00100000
	Const1Value1Mark                           = 0b00011111
	LtwValidFlagMark                           = 0b10000000
	LtwOffsetMark                              = 0b01111111
	PiecewiseRateReservedMark                  = 0b11000000
	PiecewiseRateMark                          = 0b00111111
	SpliceTypeMark                             = 0b11110000
	DtsNextAu0Mark                             = 0b00001110
	MarkerBit0Mark                             = 0b00000001
	MarkerBit1Mark                             = 0b00000001
	MarkerBit2Mark                             = 0b00000001
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
	endIndex := index + int(adaptationFieldLength)

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
		if len(buf[index:]) < 6 {
			return 0, fmt.Errorf("pcr: len(buf[index:]) < 6")
		}
		pcrbBuf := make([]byte, 5)
		pcrbBuf[0] = buf[index] >> 7
		pcrbBuf[1] = buf[index]<<1 + buf[index+1]>>7
		pcrbBuf[2] = buf[index+1]<<1 + buf[index+2]>>7
		pcrbBuf[3] = buf[index+2]<<1 + buf[index+3]>>7
		pcrbBuf[4] = buf[index+3]<<1 + buf[index+4]>>7
		programClockReferenceBase, err := util.BytesToUint64ByBigEndian(pcrbBuf)
		if err != nil {
			return 0, fmt.Errorf("util.BytesToUint64ByBigEndian(pcrbBuf) failed, err:%v", err)
		}
		fmt.Printf("programClockReferenceBase is %v\n", programClockReferenceBase)

		index += 4

		const1Value0 := (buf[index] & Const1Value0Mark) >> 1
		if const1Value0 != 0b00111111 {
			return 0, fmt.Errorf("const1Value0 != 0b00111111")
		}
		fmt.Printf("const1Value0 is 0b00111111\n")

		pcreBuf := make([]byte, 2)
		pcreBuf[0] = buf[index] & ProgramClockReferenceExtensionMark
		pcreBuf[1] = buf[index+1]

		programClockReferenceExtension, err := util.BytesToUint16ByBigEndian(pcreBuf)
		if err != nil {
			return 0, fmt.Errorf("util.BytesToUint16ByBigEndian(pcreBuf) failed, err:%v", err)
		}
		fmt.Printf("programClockReferenceExtension is %v\n", programClockReferenceExtension)

		index += 2
	}

	if opcrFlag == 1 {
		if len(buf[index:]) < 6 {
			return 0, fmt.Errorf("pcr: len(buf[index:]) < 6")
		}
		opcrbBuf := make([]byte, 5)
		opcrbBuf[0] = buf[index] >> 7
		opcrbBuf[1] = buf[index]<<1 + buf[index+1]>>7
		opcrbBuf[2] = buf[index+1]<<1 + buf[index+2]>>7
		opcrbBuf[3] = buf[index+2]<<1 + buf[index+3]>>7
		opcrbBuf[4] = buf[index+3]<<1 + buf[index+4]>>7
		originalProgramClockReferenceBase, err := util.BytesToUint64ByBigEndian(opcrbBuf)
		if err != nil {
			return 0, fmt.Errorf("util.BytesToUint64ByBigEndian(opcrbBuf) failed, err:%v", err)
		}
		fmt.Printf("originalProgramClockReferenceBase is %v\n", originalProgramClockReferenceBase)
		index += 4

		const1Value2 := (buf[index] & Const1Value2Mark) >> 1
		fmt.Printf("const1Value2 is %v\n", const1Value2)

		opcreBuf := make([]byte, 2)
		opcreBuf[0] = buf[index] & OriginalProgramClockReferenceExtensionMark
		opcreBuf[1] = buf[index+1]
		originalProgramClockReferenceExtension, err := util.BytesToUint16ByBigEndian(opcreBuf)
		if err != nil {
			return 0, fmt.Errorf("originalProgramClockReferenceExtension is %v\n", originalProgramClockReferenceExtension)
		}
		fmt.Printf("originalProgramClockReferenceExtension is %v\n", originalProgramClockReferenceExtension)

		index += 2
	}

	if splicingPointFlag == 1 {
		if len(buf[index:]) < 1 {
			return 0, fmt.Errorf("splicingPointFlag: len(buf[index:]) < 1")
		}
		spliceCountDown := util.BytesToUint8ByBigEndian(buf[index])
		fmt.Printf("spliceCountDown is %v\n", spliceCountDown)
		index += 1
	}

	if transportPrivateDataFlag == 1 {
		if len(buf[index:]) < 1 {
			return 0, fmt.Errorf("transportPrivateDataFlag: len(buf[index:]) < 1 ")
		}
		transportPrivateDataLength := util.BytesToUint8ByBigEndian(buf[index])
		fmt.Printf("transportPrivateDataLength is %v\n", transportPrivateDataLength)

		index += 1

		if len(buf[index:]) < int(transportPrivateDataLength) {
			return 0, fmt.Errorf("len(buf[index:]) < int(transportPrivateDataLength)")
		}

		transportPrivateData := buf[index : index+int(transportPrivateDataLength)/8]
		fmt.Printf("transportPrivateData is %v", transportPrivateData)

		index += int(transportPrivateDataLength) / 8
	}

	if adaptationFieldExtensionFlag == 1 {
		if len(buf[index:]) > 2 {
			return 0, fmt.Errorf("adaptationFieldExtensionFlag: len(buf[index:]) > 2")
		}
		adaptationFieldExtensionLength := util.BytesToUint8ByBigEndian(buf[index])
		fmt.Printf("adaptationFieldExtensionLength is %v\n", adaptationFieldExtensionLength)

		index += 1

		ltwFlag := (buf[index] & LtwFlagMark) >> 7
		fmt.Printf("ltwFlag is %v\n", ltwFlag)

		piecewiseRateFlag := (buf[index] & PiecewiseRateFlagMark) >> 6
		fmt.Printf("piecewiseRateFlag is %v\n", piecewiseRateFlag)

		seamlessSpliceFlag := (buf[index] & SeamlessSpliceFlagMark) >> 5
		fmt.Printf("seamlessSpliceFlag is %v\n", seamlessSpliceFlag)

		const1Value1 := buf[index] & Const1Value1Mark
		fmt.Printf("const1Value1 is %v\n", const1Value1)

		index += 1

		if ltwFlag == 1 {
			if len(buf[index:]) < 2 {
				return 0, fmt.Errorf("ltw: len(buf[index:]) < 2")
			}

			ltwValidFlag := (buf[index] & LtwValidFlagMark) >> 7
			fmt.Printf("ltwValidFlag is %v\n", ltwValidFlag)

			ltwOffset := buf[index : index+2]
			ltwOffset[0] = ltwOffset[0] & LtwOffsetMark
			fmt.Printf("ltwOffset is %v\n", ltwOffset)

			index += 2
		}

		if piecewiseRateFlag == 1 {
			if len(buf[index:]) < 3 {
				return 0, fmt.Errorf("piecewiseRate: len(buf[index:]) < 3")
			}

			piecewiseRateReserved := (buf[index] & PiecewiseRateReservedMark) >> 6
			fmt.Printf("piecewiseRateReserved is %v\n", piecewiseRateReserved)

			piecewiseRate := buf[index : index+3]
			piecewiseRate[0] = piecewiseRate[0] & PiecewiseRateMark
			fmt.Printf("piecewiseRate is %v\n", piecewiseRate)

			index += 3
		}

		if seamlessSpliceFlag == 1 {

			if len(buf[index:]) < 5 {
				return 0, fmt.Errorf("seamlessSplice: len(buf[index:]) < 5")
			}

			spliceType := (buf[index] & SpliceTypeMark) >> 4
			fmt.Printf("spliceType is %v\n", spliceType)

			dtsNextAu0 := (buf[index] & DtsNextAu0Mark) >> 1
			fmt.Printf("dtsNextAu0 is %v\n", dtsNextAu0)

			markerBit0 := buf[index] & MarkerBit0Mark
			fmt.Printf("markerBit0 is %v\n", markerBit0)

			index += 1

			dtsNextAu1 := buf[index : index+2]
			dtsNextAu1[1] = dtsNextAu1[0]<<7 + dtsNextAu1[1]>>1
			dtsNextAu1[0] = dtsNextAu1[0] >> 1
			fmt.Printf("dtsNextAu1 is %v\n", dtsNextAu1)

			index += 1

			markerBit1 := buf[index] & MarkerBit1Mark
			fmt.Printf("markerBit1 is %v\n", markerBit1)

			index += 1

			dtsNextAu2 := buf[index : index+2]
			dtsNextAu2[1] = dtsNextAu2[0]<<7 + dtsNextAu2[1]>>1
			dtsNextAu2[0] = dtsNextAu2[0] >> 1
			fmt.Printf("dtsNextAu2 is %v\n", dtsNextAu2)

			index += 1

			markerBit2 := buf[index] & MarkerBit2Mark
			fmt.Printf("markerBit2 is %v\n", markerBit2)

			index += 1
		}
	}

	fmt.Println()

	return endIndex, nil
}
