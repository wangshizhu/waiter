package main

import (
	"waiter/log"

	"github.com/rs/zerolog"
)

func main() {
	log.Init()
	log.AddHook(zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		levelName := level.String()
		if level == zerolog.NoLevel {
			levelName = "nolevel"
		}
		e.Str("level_name", levelName)
	}))

	log.Info().Msg("test hook")

	select {}
}
