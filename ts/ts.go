package ts

import "fmt"

type Ts struct {
	PacketNum                   int
	CurrentTransportStream      *TransportStream
	ProgramAssociationTable     *ProgramAssociationTable
	ProgramMapTable             *ProgramMapTable
	PacketizedElementaryStreams *PacketizedElementaryStreams
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

		fmt.Printf("----------------------------------------\n")
		fmt.Printf("ts packet: %v\n", t.PacketNum)

		t.PacketNum++

		tsBuf := buf[index : index+TransportStreamLength]
		fmt.Printf("tsBuf:\n%x\n", tsBuf)
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
		t.ProgramAssociationTable = BuildProgramAssociationTable(t)
	}

	return t.ProgramAssociationTable
}

func (t *Ts) getProgramMapTable() *ProgramMapTable {
	if t.ProgramMapTable == nil {
		t.ProgramMapTable = BuildProgramMapTable(t)
	}

	return t.ProgramMapTable
}

func (t *Ts) getPacketizedElementaryStreams() *PacketizedElementaryStreams {
	if t.PacketizedElementaryStreams == nil {
		t.PacketizedElementaryStreams = BuildPacketizedElementaryStreams(t)
	}

	return t.PacketizedElementaryStreams
}
