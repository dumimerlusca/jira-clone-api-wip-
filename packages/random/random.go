package random

import (
	"math/rand"
	"strings"
	"time"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func RandomInt(min, max int64) int64 {
	return min + seededRand.Int63n(max-min+1)
}

const alphabet = "abcdefghigklmnopqrstuvwxyz"

func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
