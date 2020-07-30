package ts

import (
	"fmt"
)

type Payload struct {
	TransportStream *TransportStream

	HavePayload bool
}

func BuildPayload(t *TransportStream) *Payload {

	p := new(Payload)
	p.TransportStream = t

	return p
}

func (p *Payload) Parse(buf []byte, index int) (int, error) {

	adaptationFieldControl := p.TransportStream.TsHeader.AdaptationFieldControl
	if adaptationFieldControl != AdaptationFieldControl01 &&
		adaptationFieldControl != AdaptationFieldControl11 {
		return index, nil
	}
	p.HavePayload = true

	var err error
	if p.TransportStream.TsHeader.Pid == 0 {
		programAssociationTable := p.TransportStream.Ts.getProgramAssociationTable()
		index, err = programAssociationTable.Parse(buf, index)
		if err != nil {
			return 0, fmt.Errorf("programAssociationTable.Parse failed, err:%v", err)
		}
	} else if p.TransportStream.Ts.ProgramAssociationTable != nil {
		if p.TransportStream.TsHeader.Pid == p.TransportStream.Ts.ProgramAssociationTable.ProgramMapPid {
			programMapTable := p.TransportStream.Ts.getProgramMapTable()
			index, err = programMapTable.Parse(buf, index)
			if err != nil {
				return 0, fmt.Errorf("programMapTable.Parse failed, err:%v", err)
			}
		} else {
			_, ok := p.TransportStream.Ts.ProgramMapTable.PidStreamTypeMap[p.TransportStream.TsHeader.Pid]
			if !ok {
				if p.TransportStream.TsHeader.Pid == 0x11 {
					fmt.Println("Pid == 0x11, jumping")
					return len(buf), nil
				}
				return 0, fmt.Errorf("p.TransportStream.Ts.ProgramMapTable.PidStreamTypeMap[p.TransportStream.TsHeader.Pid] failed")
			}
			packetizedElementaryStreams := p.TransportStream.Ts.getPacketizedElementaryStreams()
			index, err = packetizedElementaryStreams.Parse(buf, index)
			if err != nil {
				return 0, fmt.Errorf("packetizedElementaryStreams.Parse(buf, index) failed, err:%v", err)
			}
		}
	}

	return index, nil
}
