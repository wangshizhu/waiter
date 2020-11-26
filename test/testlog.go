package main

import (
	"fmt"
	"sync"
	"time"
	"waiter/log"
)

var wg sync.WaitGroup

func PrintLog() {
	defer wg.Done()
	// time.Sleep(time.Second)

	start := time.Now().UnixNano() / 1e6

	for i := 0; i < 1000; i++ {
		log.Info().Int("test", i).Msg("print log")
	}

	end := time.Now().UnixNano() / 1e6

	fmt.Println("one goroutine total:", end-start)
}

func main() {
	log.Init()
	log.DisableStdOut()

	wg.Add(100)

	start := time.Now().UnixNano() / 1e6

	for i := 0; i < 100; i++ {
		go PrintLog()
	}

	wg.Wait()

	end := time.Now().UnixNano() / 1e6

	fmt.Println("total:", end-start)

	// log.Info().Int("test", 10).Msg("init success")
	select {}
}
