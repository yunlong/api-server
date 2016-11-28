package utils

import (
	"fmt"
	"strings"
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

func LogFilenameGen() string {
	newName := RandStringRunes(10)
	return newName + ".txt"
}

func UuidGenerated() string {
	u1 := uuid.NewV4()
	return fmt.Sprintf("%s", u1)
}

func RenameImage(oldName string, nameLength int) string {
	imageEtx := oldName[strings.LastIndex(oldName,"."):len(oldName)]
	newName := RandStringRunes(nameLength)
	return newName + imageEtx
}

func AddProjectEnv(m map[string]string, key, value string) {
    _, madd := m[key]
    if !madd {
		m[key] = value
    }
}
