package ts

import "fmt"

type Ts struct {
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
		ts := new(TransportStream)
		err := ts.Parse(tsBuf)
		if err != nil {
			return nil, fmt.Errorf("ts.Parse(tsBuf) failed, err:%v", err)
		}
		index += TransportStreamLength
	}

	return buf[index:], nil
}
