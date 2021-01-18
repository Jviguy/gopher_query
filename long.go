/**

        :::   :::       :::     :::::::::  ::::::::::
      :+:+: :+:+:    :+: :+:   :+:    :+: :+:
    +:+ +:+:+ +:+  +:+   +:+  +:+    +:+ +:+
   +#+  +:+  +#+ +#++:++#++: +#+    +:+ +#++:++#
  +#+       +#+ +#+     +#+ +#+    +#+ +#+
 #+#       #+# #+#     #+# #+#    #+# #+#
###       ### ###     ### #########  ##########

            :::::::::  :::   :::
            :+:    :+: :+:   :+:
            +:+    +:+  +:+ +:+
            +#++:++#+    +#++:
            +#+    +#+    +#+
            #+#    #+#    #+#
            #########     ###

     ::::::::::: :::     ::: ::::::::::: ::::::::  :::    ::: :::   :::
        :+:     :+:     :+:     :+:    :+:    :+: :+:    :+: :+:   :+:
       +:+     +:+     +:+     +:+    +:+        +:+    +:+  +:+ +:+
      +#+     +#+     +:+     +#+    :#:        +#+    +:+   +#++:
     +#+      +#+   +#+      +#+    +#+   +#+# +#+    +#+    +#+
#+# #+#       #+#+#+#       #+#    #+#    #+# #+#    #+#    #+#
 #####          ###     ########### ########   ########     ###

MIT License

Copyright (c) 2020 Jviguy

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package gopher_query

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
)

const (
	Magik = 0xFEFD
	Stat = 0x00
	Handshake = 0x09
)

type LongQueryResponse struct {
	ServerSoftware string
	Plugins string
	Version string
	Whitelist string
	Players []string
	PlayerCount string
	MaxPlayers string
	GameName string
	GameMode string
	MapName string
	HostName string
	HostIp string
	HostPort string
}

func (c Client) LongQuery(addr string) (LongQueryResponse, error) {
	sid := rand.Int31()
	resp := LongQueryResponse{}
	conn, err := c.dialer.Dial("udp", addr)
	if err != nil {
		return resp, err
	}
	ct, err := generateChallengeToken(conn, sid)
	if err != nil {
		return resp, err
	}
	data, players, err := fullStat(conn, sid, ct)
	if err != nil {
		return resp, err
	}
	resp.Players = players
	resp.ServerSoftware = data["server_engine"]
	resp.Plugins = data["plugins"]
	resp.Whitelist = data["whitelist"]
	resp.Version = data["version"]
	resp.PlayerCount = data["numplayers"]
	resp.MaxPlayers = data["maxplayers"]
	resp.MapName = data["map"]
	resp.HostPort = data["hostport"]
	resp.HostName = data["hostname"]
	resp.HostIp = data["hostip"]
	resp.GameMode = data["gametype"]
	resp.GameName = data["game_id"]
	return resp, err
}

func fullStat(conn net.Conn, sid int32, ct int32) (map[string]string, []string, error) {
	buf := &strings.Builder{}
	info := make(map[string]string)
	err := binary.Write(buf, binary.BigEndian, uint16(Magik))
	if err != nil {
		return info, nil, err
	}
	buf.WriteByte(Stat)
	err = binary.Write(buf, binary.BigEndian, sid)
	if err != nil {
		return info, nil, err
	}
	err = binary.Write(buf, binary.BigEndian, ct)
	if err != nil {
		return info, nil, err
	}
	//Pad 4
	buf.Write([]byte{0, 0, 0, 0})
	_, err = conn.Write([]byte(buf.String()))
	if err != nil {
		return info, nil, err
	}
	buf.Reset()
	tmp := make([]uint8, math.MaxUint16)
	_, err = conn.Read(tmp)
	tmp = bytes.TrimRight(tmp, "\x00")
	if err != nil {
		return info, nil, err
	}
	id := tmp[0]
	if id == Stat {
		playerKey := [...]byte{0x00, 0x01, 'p', 'l', 'a', 'y', 'e', 'r', '_', 0x00, 0x00}
		bs := tmp[16:]
		data := bs
		playerIndex := bytes.Index(bs, playerKey[:])
		if playerIndex != -1 {
			bs = bs[:playerIndex]
		}
		var wg sync.WaitGroup
		go func() {
			wg.Add(1)
			vals := bytes.Split(bs, []byte{0x00})

			if len(vals) % 2 != 0 {
				vals = vals[:len(vals)-1]
			}
			for i := 0; i < len(vals); i += 2 {
				info[string(vals[i])] = string(vals[i+1])
			}
			wg.Done()
		}()
		if playerIndex != -1 {
			pD := data[playerIndex+len(playerKey):]
			vals := bytes.Split(pD, []byte{0x00})
			players := make([]string, 0, len(vals))
			go func() {
				wg.Add(1)
				for i := 0; i < len(vals); i++ {
					if len(vals[i]) == 0 {
						break
					}
					players = append(players, string(vals[i]))
				}
				wg.Done()
			}()
			wg.Wait()
			return info, players, nil
		}
		return info, nil, nil
	}
	return info, nil, fmt.Errorf("invalid packet recieved while awaiting resp, id: %v", id)
}

/**
Used to generate a challenge token for a full stat query.
 */
func generateChallengeToken(conn net.Conn, sid int32) (int32, error) {
	buf := &strings.Builder{}
	err := binary.Write(buf, binary.BigEndian, uint16(Magik))
	if err != nil {
		return 0, err
	}
	buf.WriteByte(Handshake)
	err = binary.Write(buf, binary.BigEndian, sid)
	if err != nil {
		return 0, err
	}
	_, err = conn.Write([]byte(buf.String()))
	if err != nil {
		return 0, err
	}
	buf.Reset()
	tmp := make([]byte, math.MaxUint16 >> 2)
	_, err = conn.Read(tmp)
	if err != nil {
		return 0, err
	}
	buf.Write(bytes.Trim(tmp, "\x00"))
	id := []byte(buf.String())[0]
	if id == Handshake {
		str := buf.String()[binary.Size(sid)+1:]
		ct, err := strconv.ParseInt(str, 0, 32)
		if err != nil {
			return 0, err
		}
		return int32(ct), err
	}
	return 0, fmt.Errorf("invalid packet recieved while awaiting handshake resp, id: %v", id)
}