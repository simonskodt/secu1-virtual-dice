package utilities

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"runtime"
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

	i := 1
	for i < n {
		strBuilder.WriteString(strconv.Itoa(rand.Intn(2)))
		i++
	}

	return strBuilder.String()
}

func ConcatStrings(this string, other string) string {
	return this + other
}

func HashStr(str string) string {
	hashCode := sha256.New()
	hashCode.Write([]byte(str))

	return hex.EncodeToString(hashCode.Sum(nil))
}

func ExclusiveOrOnTwoDiceResultsMod6Plus1(this int, other int) int {
	xor := this ^ other
	return (xor % 6) + 1
}

var Reset = "\033[0m"
var Green  = "\033[32m"

func SetColorOfLine() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Green = ""
	}
}