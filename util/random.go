package util

import (
	"math/rand"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt returns, as an int64, a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}

// RandomDecimal returns, as a decimal.Decimal, a random decimal between min and max
func RandomDecimal(min, max int64) decimal.Decimal {
	return decimal.NewFromInt(min).Add(decimal.NewFromInt(rand.Int63n(max - min + 1)))
}

// RandomString returns a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}