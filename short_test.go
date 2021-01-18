package gopher_query

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestClient_ShortQuery(t *testing.T) {
	resp := ShortQueryResponse{}
	conn, err := net.Dial("udp", "127.0.0.1:19133")
	if err != nil {
		t.Fatal(err)
	}
	magik, err := hex.DecodeString("00ffff00fefefefefdfdfdfd12345678")
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	buf.WriteByte(0x01)
	err = binary.Write(&buf, binary.BigEndian, time.Now().Unix())
	if err != nil {
		t.Fatal(err)
	}
	buf.Write(magik)
	err = binary.Write(&buf, binary.BigEndian, rand.Int63())
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	buf.Reset()
	var tmp = make([]uint8, math.MaxUint16)
	_,err = conn.Read(tmp)
	if err != nil {
		t.Fatal(err)
	}
	body := strings.Split(string(tmp[16+len(magik):]), ";")
	resp.GameEdition = body[0]
	resp.MOTD = make([]string, 2)
	resp.MOTD[0] = body[1]
	resp.MOTD[1] = body[7]
	proto, err := strconv.Atoi(body[2])
	if err != nil {
		t.Fatal(body)
	}
	resp.ProtocolVersion = proto
	resp.GameVersion = body[3]
	pc, err := strconv.Atoi(body[4])
	if err != nil {
		t.Fatal(err)
	}
	resp.PlayerCount = pc
	mpc, err := strconv.Atoi(body[5])
	if err != nil {
		t.Fatal(err)
	}
	resp.MaxPlayerCount = mpc
	resp.ServerUID = body[6]
	resp.GameMode = body[8]
	t.Log(resp)
}
