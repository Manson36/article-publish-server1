package main

import (
	"crypto/rand"
	"fmt"
)

func main() {

	var seedString = "2345678abcdefhijkmnprstwxyzABCDEFGHJKMNPQRSTWXYZ"
	x := make([]byte, 4)
	out := make([]byte, 4)

	_, err := rand.Read(x)
	if err != nil {
		panic(err)
	}
	for id, key := range x {
		m := byte(int(key) % len(seedString))
		out[id] = seedString[m]
	}

	fmt.Println(string(out))
}
