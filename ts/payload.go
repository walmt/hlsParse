package ts

import (
	"fmt"
)

type Payload struct {
	TransportStream         *TransportStream
}

func (p *Payload) Parse(buf []byte, index int) (int, error) {
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
		}
	}

	return 0, nil
}

