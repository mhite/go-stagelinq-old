package stagelinq

import "net"

// BeatInfo represents a received beat value.
type BeatInfo struct {
	Name  string
	Value map[string]interface{}
}

// BeatInfoConnection provides functionality to communicate with the BeatInfo data source.
type BeatInfoConnection struct {
	conn      *messageConnection
	errC      chan error
	beatInfoC chan *BeatInfo
}

var beatInfoConnectionMessageSet = newDeviceConnMessageSet([]message{&beatEmitMessage{}})

func NewBeatInfoConnection(conn net.Conn, token Token) (bic *BeatInfoConnection, err error) {
	msgConn := newMessageConnection(conn, beatInfoConnectionMessageSet)
	return &BeatInfoConnection{
		conn:      msgConn,
		errC:      make(chan error),
		beatInfoC: make(chan *BeatInfo),
	}, err
}

// Subscribe tells the StagelinQ device to let us know about changes for the given state value.
func (bic *BeatInfoConnection) Subscribe() error {
	return bic.conn.WriteMessage(&beatInfoSubscribeMessage{})
}
