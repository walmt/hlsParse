package ts

import (
	"fmt"
)

const (
	//Pes10Mark                 = 0b11000000
	//PesScramblingControl      = 0b00110000
	//PesPriority               = 0b00001000
	//PesDataAlignmentIndicator = 0b00000100
	//PesCopyRight              = 0b00000010
	//PesOriginalOrCopy         = 0b00000001
	//PesPtsDtsFlags            = 0b11000000
	//PesEscrFlag               = 0b00100000
	//PesEsRateFlag             = 0b00010000
	//PesDsmTrickModeFlag       = 0b00001000
	//PesAdditionalCopyInfoFlag = 0b00000100
	//PesCrcFlag                = 0b00000010
	//PesExtensionFlag          = 0b00000001

	PtsConstValue0Mark = 0b11110000
	PtsPtsAMark        = 0b00001110
	PtsConstValue1Mark = 0b00000001
	PtsPtsBMark        = 0b11111110
	PtsConstValue2Mark = 0b00000001
	PtsPtsCMark        = 0b11111110
	PtsConstValue3Mark = 0b00000001

	DtsConstValue0Mark = 0b11110000
	DtsDtsAMark        = 0b00001110
	DtsConstValue1Mark = 0b00000001
	DtsDtsBMark        = 0b11111110
	DtsConstValue2Mark = 0b00000001
	DtsDtsCMark        = 0b11111110
	DtsConstValue3Mark = 0b00000001
)

type PacketizedElementaryStreams struct {
	Ts                  *Ts
	ElementaryStreamMap map[uint16]*ElementaryStreams
}

func BuildPacketizedElementaryStreams(t *Ts) *PacketizedElementaryStreams {

	p := new(PacketizedElementaryStreams)
	p.ElementaryStreamMap = make(map[uint16]*ElementaryStreams)
	p.Ts = t

	return p
}

func (p *PacketizedElementaryStreams) Parse(buf []byte, index int) (int, error) {

	var err error
	e, ok := p.ElementaryStreamMap[p.Ts.CurrentTransportStream.TsHeader.Pid]
	if !ok {
		e = buildElementaryStream(p)
		p.ElementaryStreamMap[p.Ts.CurrentTransportStream.TsHeader.Pid] = e
		index, err = e.ParseHeader(buf, index)
		if err != nil {
			return 0, fmt.Errorf("e.ParseHeader(buf, index) failed, err:%v", err)
		}
	}
	index, err = e.ParseBody(buf, index)
	if err != nil {
		return 0, fmt.Errorf("e.ParseBody(buf, index) failed, err:%v", err)
	}

	fmt.Println()
	
	return index, nil
}
