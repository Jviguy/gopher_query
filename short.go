package gopher_query

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type ShortQueryResponse struct {
	GameEdition string
	//A slice to represent the two given lines of motd there are.
	MOTD            []string
	ProtocolVersion int
	GameVersion     string
	PlayerCount     int
	MaxPlayerCount  int
	ServerUID       string
	GameMode        string
	GameModeInteger int
	Port            uint16
	PortV6          uint16
}

func (c Client) ShortQuery(addr string) (ShortQueryResponse, error) {
	resp := ShortQueryResponse{}
	conn, err := c.dialer.Dial("raknet", addr)
	if err != nil {
		return resp, err
	}
	magik, err := hex.DecodeString("00ffff00fefefefefdfdfdfd12345678")
	if err != nil {
		return resp, err
	}
	var buf bytes.Buffer
	buf.WriteByte(0x01)
	err = binary.Write(&buf, binary.BigEndian, time.Now().Unix())
	if err != nil {
		return resp, err
	}
	buf.Write(magik)
	err = binary.Write(&buf, binary.BigEndian, rand.Int63())
	if err != nil {
		return resp, err
	}
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return resp, err
	}
	buf.Reset()
	var tmp = make([]uint8, math.MaxUint16)
	_, err = conn.Read(tmp)
	if err != nil {
		return resp, err
	}
	body := strings.Split(string(tmp), ";")
	resp.GameEdition = body[0]
	resp.MOTD = make([]string, 2)
	resp.MOTD[0] = body[1]
	resp.MOTD[1] = body[7]
	proto, err := strconv.Atoi(body[2])
	if err != nil {
		return resp, err
	}
	resp.ProtocolVersion = proto
	resp.GameVersion = body[3]
	pc, err := strconv.Atoi(body[4])
	if err != nil {
		return resp, err
	}
	resp.PlayerCount = pc
	mpc, err := strconv.Atoi(body[5])
	if err != nil {
		return resp, err
	}
	resp.MaxPlayerCount = mpc
	resp.ServerUID = body[6]
	resp.GameMode = body[8]
	return resp, err
}
