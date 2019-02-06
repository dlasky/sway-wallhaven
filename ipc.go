package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"unsafe"
)

var nativeEndian binary.ByteOrder

func init() {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		nativeEndian = binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		nativeEndian = binary.BigEndian
	default:
		panic("Could not determine native endianness.")
	}
}

type messageType uint32

const (
	messageTypeRunCommand messageType = iota
	messageTypeGetWorkspaces
	messageTypeSubscribe
	messageTypeGetOutputs
	messageTypeGetTree
	messageTypeGetMarks
	messageTypeGetBarConfig
	messageTypeGetVersion
	messageTypeGetBindingModes
	messageTypeGetConfig
	messageTypeSendTick
)

const (
	messageReplyTypeCommand messageType = iota
	messageReplyTypeWorkspaces
	messageReplyTypeSubscribe
)

var magic = [6]byte{'i', '3', '-', 'i', 'p', 'c'}

type header struct {
	Magic  [6]byte
	Length uint32
	Type   messageType
}

type message struct {
	Type    messageType
	Payload []byte
}

func getSwayIPCPath() string {
	path := os.Getenv("SWAYSOCK")
	if path == "" {
		out, err := exec.Command("sway", "--get-socketpath").CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
		path = strings.TrimSpace(string(out))
	}
	return path

}

func getSocket() (net.Conn, error) {
	conn, err := net.Dial("unix", getSwayIPCPath())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func write(conn net.Conn, msg message) error {
	if err := binary.Write(conn, nativeEndian, &header{magic, uint32(len(msg.Payload)), msg.Type}); err != nil {
		return err
	}
	if len(msg.Payload) > 0 {
		_, err := conn.Write(msg.Payload)
		if err != nil {
			return err
		}
	}
	return nil
}

func read(conn net.Conn) (message, error) {
	var head header
	if err := binary.Read(conn, nativeEndian, &head); err != nil {
		return message{}, err
	}
	msg := message{
		Type:    head.Type,
		Payload: make([]byte, head.Length),
	}
	_, err := io.ReadFull(conn, msg.Payload)
	return msg, err
}

func trip(conn net.Conn, msg message) (message, error) {
	err := write(conn, msg)
	if err != nil {
		return message{}, err
	}
	return read(conn)
}
