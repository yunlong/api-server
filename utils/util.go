package utils

import (
	//"crypto/rand"
	"fmt"
	"math/rand"
	//"log"
	//"time"

	"github.com/satori/go.uuid"

	//"github.com/heroku/docker-registry-client/registry"
	"strings"
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

/*
func GetTagLatest(image string) (string, error) {
	reg := registry.TagsFullResponse{}
	url := "https://hub.docker.com/"
	username := "" // anonymous
	password := "" // anonymous
	if hub, err := registry.New(url, username, password); err != nil {
		log.Printf("failed at creating Hub client")
		return "", err
	} else if tags, err := hub.FullTags(image); err != nil {
		log.Printf("failed at getting tag metadata")
		return "", err
	} else {
		newest := time.Time{}
		tagLatest := ""
		//layout := "2000-01-01T21:06:44.982740Z"
		for _, result := range tags.Results {
			t, err := time.Parse(time.RFC3339, result.LastUpdated)
			if err != nil {
				log.Fatal(err)
			}

			deltaT := time.Now().Sub(t)
			deltaN := time.Now().Sub(newest)
			if deltaT < deltaN {
				newest = t
				tagLatest = result.Name
			}
		}
		return tagLatest, nil
	}
}*/
