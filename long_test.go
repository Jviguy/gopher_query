package gopher_query

import (
	"math/rand"
	"net"
	"testing"
)

func TestClient_LongQuery(t *testing.T) {
	sid := rand.Int31()
	conn, err := net.Dial("udp", "versai.pro:19132")
	if err != nil {
		t.Fatal(err)
	}
	ct, err := generateChallengeToken(conn, sid)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ct)
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