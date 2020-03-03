package utils

import (
	"crypto/rand"
	"github.com/article-publish-server1/datamodels"
	"strconv"
)

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

func StringSliceToJsonNumArray(arr []string) (datamodels.JsonNumArray, error) {
	if len(arr) == 0 {
		return nil, nil
	}

	numArr := make(datamodels.JsonNumArray, 0)
	for _, v := range arr {
		num, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}

		numArr = append(numArr, num)
	}

	return numArr, nil
}
