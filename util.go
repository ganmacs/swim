package swim

import (
	"math/rand"
	"net"
	"strconv"
)

func joinHostPort(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}

func randInt(n int) int {
	if n == 0 {
		return n
	}

	return (rand.Int() % n)
}
