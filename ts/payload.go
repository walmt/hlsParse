package ts

import "fmt"

type Payload struct {
	TransportStream         *TransportStream
	ProgramAssociationTable *ProgramAssociationTable
}

func (p *Payload) Parse(buf []byte, index int) (int, error) {
	var err error
	if p.TransportStream.TsHeader.Pid == 0 {
		programAssociationTable := p.getProgramAssociationTable()
		index, err = programAssociationTable.Parse(buf, index)
		if err != nil {
			return 0, fmt.Errorf("programAssociationTable.Parse failed, err:%v", err)
		}
	}

	return 0, nil
}

func (p *Payload) getProgramAssociationTable() *ProgramAssociationTable {
	if p.ProgramAssociationTable == nil {
		p.ProgramAssociationTable = new(ProgramAssociationTable)
		p.ProgramAssociationTable.Payload = p
	}

	return p.ProgramAssociationTable
}
