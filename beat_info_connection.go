package stagelinq

import (
	"net"
)

// BeatInfo represents a received BeatInfo message.
type BeatInfo struct {
	Clock     uint64
	Players   []PlayerInfo
	Timelines []uint64
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

	errC := make(chan error, 1)
	beatInfoC := make(chan *BeatInfo, 1)

	beatInfoConn := BeatInfoConnection{
		conn:      msgConn,
		errC:      errC,
		beatInfoC: beatInfoC,
	}
	// perform in-protocol service request
	msgConn.WriteMessage(&serviceAnnouncementMessage{
		tokenPrefixedMessage: tokenPrefixedMessage{
			Token: token,
		},
		Service: "BeatInfo",
		Port:    uint16(getPort(conn.LocalAddr())),
	})

	go func() {
		var err error
		defer func() {
			if err != nil {
				beatInfoConn.errC <- err
				close(beatInfoConn.errC)
			}
			close(beatInfoConn.beatInfoC)
		}()
		for {
			var msg message
			msg, err = msgConn.ReadMessage()
			if err != nil {
				return
			}

			switch v := msg.(type) {
			case *beatEmitMessage:
				beatInfo := &BeatInfo{
					Clock:     v.Clock,
					Players:   v.Players,
					Timelines: v.Timelines,
				}
				beatInfoC <- beatInfo
			}
		}
	}()

	bic = &beatInfoConn
	return
}

// Subscribe tells the StagelinQ device to let us know about changes for the given state value.
func (bic *BeatInfoConnection) Subscribe() error {
	return bic.conn.WriteMessage(&beatInfoSubscribeMessage{})
}

// StateC returns the channel via which state changes will be returned for this connection.
func (bic *BeatInfoConnection) BeatInfoC() <-chan *BeatInfo {
	return bic.beatInfoC
}

// ErrorC returns the channel via which connectionrerors will be returned for this connection.
func (bic *BeatInfoConnection) ErrorC() <-chan error {
	return bic.errC
}
