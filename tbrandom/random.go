package tbrandom

import (
	"math/rand"
)

const (
	CHARS string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-+="
	CHLEN        = len(CHARS)
)

func GenRandNum(low, top int) int {
	return (rand.Intn(top-low) + low)
}

func GenRandString(length int) string {
	var result []rune
	for i := 0; i < length; i++ {
		j := GenRandNum(0, CHLEN-1)
		result = append(result, rune(CHARS[j]))
	}
	return string(result)
}
