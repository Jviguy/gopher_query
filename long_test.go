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
	"fmt"
	"math/rand"
	"net"
	"testing"
)

func TestClient_LongQuery(t *testing.T) {
	sid := rand.Int31()
	resp := LongQueryResponse{}
	conn, err := net.Dial("udp", "velvetpractice.live:19132")
	if err != nil {
		t.Fatal(err)
	}
	ct, err := generateChallengeToken(conn, sid)
	if err != nil {
		t.Fatal(err)
	}
	data, players, err := fullStat(conn, sid, ct)
	if err != nil {
		t.Fatal(err)
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
	fmt.Println(resp)
}

func BenchmarkClient_LongQuery(b *testing.B) {
	var c = NewClient()
	for i := 0; i < b.N; i++ {
		_, err := c.LongQuery("versai.pro:19132")
		if err != nil {
			b.Fatal(err)
		}
	}
}