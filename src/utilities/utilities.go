package utilities

import (
	"math/rand"
	"strconv"
	"strings"
)

func GenerateRandDiceRoll() int {
	return rand.Intn(6) + 1
}

func FormatInt(n int) string {
	return strconv.FormatInt(int64(n), 2)
}

func GenerateRandBitStr(n int) string {
	var strBuilder strings.Builder

	for i := 0; i < n; i++ {
		bit := rand.Intn(1)
		strBuilder.WriteString(strconv.Itoa(bit))
	}

	return strBuilder.String()
}

func ConcatStrings(this string, other string) string {
	return this + other
}