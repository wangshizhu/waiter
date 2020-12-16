package main

import (
	"fmt"
	"sync"
	"time"
	"waiter/log"
)

var wg sync.WaitGroup

func PrintLog(routineIndex int) {
	defer wg.Done()
	// time.Sleep(time.Second)

	// strIndex := "named_" + strconv.Itoa(routineIndex)

	start := time.Now().UnixNano() / 1e6

	for i := 0; i < 1000; i++ {
		// log.Info().Int("test", i).Msg("print log")
		log.EasyInfo().Int("test", i).Msg("print log")
		// log.NamedInfo("named").Int("test", i).Msg("print log")
	}

	end := time.Now().UnixNano() / 1e6

	fmt.Println("one goroutine total:", end-start)
}

func main() {
	log.Init()
	log.DisableStdOut()
	log.SetLevel(0)

	wg.Add(100)

	start := time.Now().UnixNano() / 1e6

	for i := 0; i < 100; i++ {
		go PrintLog(i)
	}

	wg.Wait()

	end := time.Now().UnixNano() / 1e6

	fmt.Println("total:", end-start)

	// log.Info().Int("test", 10).Msg("init success")
	select {}
}
