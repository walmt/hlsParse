package ts

import "fmt"

type Ts struct {
	CurrentTransportStream  *TransportStream
	ProgramAssociationTable *ProgramAssociationTable
	ProgramMapTable         *ProgramMapTable
}

const (
	TransportStreamLength = 188
)

func (t *Ts) Parse(buf []byte) ([]byte, error) {

	index := 0
	for true {
		if len(buf[index:]) < 188 {
			return buf[index:], nil
		}

		tsBuf := buf[index : index+TransportStreamLength]
		ts := t.createTransportStream()
		err := ts.Parse(tsBuf)
		if err != nil {
			return nil, fmt.Errorf("ts.Parse(tsBuf) failed, err:%v", err)
		}
		index += TransportStreamLength
	}

	return buf[index:], nil
}

func (t *Ts) createTransportStream() *TransportStream {
	ts := new(TransportStream)
	ts.Ts = t
	t.CurrentTransportStream = ts
	return ts
}

func (t *Ts) getProgramAssociationTable() *ProgramAssociationTable {
	if t.ProgramAssociationTable == nil {
		t.ProgramAssociationTable = new(ProgramAssociationTable)
		t.ProgramAssociationTable.Ts = t
	}

	return t.ProgramAssociationTable
}

func (t *Ts) getProgramMapTable() *ProgramMapTable {
	if t.ProgramMapTable == nil {
		t.ProgramMapTable = new(ProgramMapTable)
		t.ProgramMapTable.Ts = t
	}

	return t.ProgramMapTable
}
