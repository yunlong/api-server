package utils

import (
	//"crypto/rand"
	"fmt"
	"math/rand"

	"github.com/satori/go.uuid"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

/*func RandomHexString(byteLength uint32) (str string, err error) {
	slice := make([]byte, byteLength)
	if _, err = rand.Read(slice); err == nil {
		str = fmt.Sprintf("%x", slice)
	}
	return
}*/

func UuidGenerated() string {
	u1 := uuid.NewV4()
	return fmt.Sprintf("%s", u1)
}
