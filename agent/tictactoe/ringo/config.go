package main

import (
	"flag"
	"strings"
)

type Config struct {
	RunSlowly bool
	RedisUrl  string
	Bot       string
	Boards    int
}

func parseArgs() Config {
	slow := flag.Bool("slow", false, "play slowly")
	redisUrl := flag.String("redis-url", "redis://localhost", "URL of redis")
	boards := flag.Int("boards", 1, "Number of tic tac toe boards to play at once")

	flag.Parse()
	return Config{
		RunSlowly: *slow,
		RedisUrl:  *redisUrl,
		Boards:    *boards,
		Bot:       strings.Join(flag.Args(), " "),
	}
}
