package main

import "waiter/log"

func main() {
	log.Init()
	log.Info().Int("test", 10).Msg("init success")
	select {}
}
