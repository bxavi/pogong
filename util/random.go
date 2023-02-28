package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}
func RandomBigInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
func RandomMoney() int64 {
	return rand.Int63n(RandomBigInt(1, 1000))
}
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(10))
}
