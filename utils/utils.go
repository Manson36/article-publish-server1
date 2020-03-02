package utils

import "crypto/rand"

func GetRandomString(size int) string {
	data := make([]byte, size)
	out := make([]byte, size)
	buffer := len(seedString)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	for id, key := range data {
		x := byte(int(key) % buffer)
		out[id] = seedString[x]
	}

	return string(out)
}
