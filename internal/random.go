package internal

import "math/rand"

const letters = "abcdefghijklmnopqrstuvwxyz"

func RandomString(length int) string {
	str := make([]rune, length)
	for i := 0; i < length; i++ {
		str[i] = rune(letters[rand.Intn(len(letters))])
	}
	return string(str)
}
