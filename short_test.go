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
	conn, err := net.Dial("udp", "velvetpractice.live:19132")
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
	_, err = conn.Read(tmp)
	if err != nil {
		t.Fatal(err)
	}
	body := strings.Split(string(tmp[len(magik) + 19:]), ";")
	t.Log(body[0])
	resp.GameEdition = body[0]
	resp.MOTD = make([]string, 2)
	resp.MOTD[0] = body[1]
	resp.MOTD[1] = body[7]
	proto, err := strconv.Atoi(body[2])
	if err != nil {
		t.Fatal(err)
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
}
