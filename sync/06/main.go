package main

import (
	"fmt"
	"sync"
)

func main() {
	type Button struct {
		Clicked *sync.Cond
	}
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(con *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			defer goroutineRunning.Done()
			con.L.Lock()
			defer con.L.Unlock()
			con.Wait()
			fn()
		}()
		goroutineRunning.Wait()
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(1)
	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})
	clickRegistered.Add(1)
	subscribe(button.Clicked, func() {
		fmt.Println("Dsiplaying annoying dialog box!")
		clickRegistered.Done()
	})
	clickRegistered.Add(1)
	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked!")
		clickRegistered.Done()
	})

	button.Clicked.Broadcast()
	clickRegistered.Wait()
}
