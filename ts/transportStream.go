package ts

import (
	"fmt"
)

type TransportStream struct {
	Ts              *Ts
	TsHeader        *TsHeader
	AdaptationField *AdaptationField
	Payload         *Payload
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

	adaptationField := t.getAdaptationField()
	index, err = adaptationField.Parse(buf, index)
	if err != nil {
		return fmt.Errorf("adaptationField.Parse failed, err:%v", err)
	}

	payload := t.getPayload()
	index, err = payload.Parse(buf, index)
	if err != nil {
		return fmt.Errorf("programAssociationTable.Parse(buf, index) failde, err:%v", err)
	}

	return nil
}

func (t *TransportStream) getTsHeader() *TsHeader {

	if t.TsHeader == nil {
		t.TsHeader = new(TsHeader)
		t.TsHeader.TransportStream = t
	}

	return t.TsHeader
}

func (t *TransportStream) getAdaptationField() *AdaptationField {

	if t.AdaptationField == nil {
		t.AdaptationField = new(AdaptationField)
		t.AdaptationField.TransportStream = t
	}

	return t.AdaptationField
}

func (t *TransportStream) getPayload() *Payload {

	if t.Payload == nil {
		t.Payload = new(Payload)
		t.Payload.TransportStream = t
	}

	return t.Payload
}
