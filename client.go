package coal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

type messageType int32

const (
	messageTypeResponse messageType = iota
	_
	messageTypeCommand
	messageTypeLogin

	headerSize     = 10
	maxPayloadSize = 1446
)

var (
	ErrUnauthenticated = errors.New("invalid password")
	ErrBadPayload      = errors.New("payload too large")
)

func Connect(address string, password string) (*RCONClient, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	rc := &RCONClient{conn: conn}

	resp, err := rc.Send(password) // the first rc.Send performs login
	if err != nil {
		return nil, err
	}

	if resp.ReqId == -1 {
		return nil, ErrUnauthenticated
	}

	return rc, nil
}

func (rc *RCONClient) Send(cmd string) (Message, error) {
	if len(cmd) > maxPayloadSize {
		return Message{}, ErrBadPayload
	}

	msg := Message{
		length:  int32(len(cmd) + headerSize),
		ReqId:   rc.reqId.Add(1),
		Content: cmd,
	}

	if msg.ReqId == 1 {
		msg.msgType = messageTypeLogin
	} else {
		msg.msgType = messageTypeCommand
	}

	rc.mu.Lock()
	_, err := rc.conn.Write(encode(msg))
	if err != nil {
		return Message{}, err
	}

	buf := make([]byte, 4110)
	n, err := rc.conn.Read(buf)
	if err != nil {
		return Message{}, err
	}
	rc.mu.Unlock()

	resp, err := decode(buf[:n])
	if err != nil {
		return Message{}, nil
	}

	return resp, err
}

func (rc *RCONClient) Close() {
	rc.conn.Close()
}

func encode(msg Message) []byte {
	length := len(msg.Content) + headerSize + 4

	buf := make([]byte, 0, length)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(length-4)) // len of remaining data
	buf = binary.LittleEndian.AppendUint32(buf, uint32(msg.ReqId))
	buf = binary.LittleEndian.AppendUint32(buf, uint32(msg.msgType))
	buf = append(buf, msg.Content...)
	buf = append(buf, 0, 0)

	return buf
}

func decode(b []byte) (Message, error) {
	reader := bytes.NewReader(b)

	var length int32
	if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
		return Message{}, err
	}

	var reqId int32
	if err := binary.Read(reader, binary.LittleEndian, &reqId); err != nil {
		return Message{}, err
	}

	var respType messageType
	if err := binary.Read(reader, binary.LittleEndian, &respType); err != nil {
		return Message{}, err
	}

	msg := Message{
		length:  length,
		ReqId:   reqId,
		msgType: respType,
	}

	payloadSize := length - headerSize
	if payloadSize > 0 {
		data := make([]byte, payloadSize)
		_, err := io.ReadFull(reader, data)
		if err != nil {
			return Message{}, err
		}
		msg.Content = string(data)
	}
	return msg, nil
}

type Message struct {
	Content string
	ReqId   int32
	msgType messageType
	length  int32
}

type RCONClient struct {
	conn  net.Conn
	reqId atomic.Int32
	mu    sync.Mutex
}
