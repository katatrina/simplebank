package util

import (
	"math/rand"
	"strings"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

// RandomInt64 returns a random number between min and max.
func RandomInt64(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString returns a random string of length n.
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner returns a random owner name with 6 characters.
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney returns a random amount of money between 0 and 1000.
func RandomMoney() int64 {
	return RandomInt64(0, 1000)
}

// RandomCurrency returns a random currency code.
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
