package common

import (
	"fmt"
	"time"
)

func TimeConsume(start time.Time) {
	fmt.Printf("cost %s\n", time.Since(start).String())
}

func IsByteSame(data1 []byte, data2 []byte) bool {
	if len(data1) != len(data2) {
		return false
	}
	for i := 0; i < len(data1); i++ {
		if data1[i] != data2[2] {
			return false
		}
	}
	return true
}
