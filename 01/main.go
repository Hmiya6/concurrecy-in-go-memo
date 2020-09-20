// This code is NOT good one!
package main

import (
	"fmt"
	"sync"
)

func main() {
	var memoryAccess sync.Mutex // (1)
	var data int
	go func() {
		memoryAccess.Lock() // (2)
		data++
		memoryAccess.Unlock() // (3)
	}()

	memoryAccess.Lock() // (4)
	if data == 0 {
		fmt.Printf("the value is %v\n", data)
	}
	memoryAccess.Unlock() // (5)
}
