package helper

import (
	"github.com/patrickmn/go-cache"
	"math/rand"
	"time"
)

var Cache = cache.New(5*time.Minute, 5*time.Minute)

type Emp []struct {
	ID        int `json:"id"`
	Name string `json:"name"`
}

func SetCache(key string, emp interface{}) bool {
	Cache.Set(key, emp, cache.NoExpiration)
	return true
}

func GetCache(key string) (string, bool) {
	var emp string
	var found bool
	data, found := Cache.Get(key)
	if found {
		emp = data.(string)
	}

	return emp, found
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789"
	bytes, err := generateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}
func GenerateRandomStringWithChar(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes, err := generateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}