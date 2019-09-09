package common

import (
	"fmt"
	"time"
)

func TimeConsume(start time.Time) {
	fmt.Printf("cost %s\n", time.Since(start).String())
}
