package main

import (
	"fmt"
	tictactoe "github.com/ai-battleground/codemelee/client/tictactoe/redis"
	"github.com/rs/xid"
	"time"
)

const logTimeFormat = "2006-01-02 15:04:05.000"

func main() {
	config := parseArgs()

	driver, err := tictactoe.NewDriver(config.RedisUrl)
	if err != nil {
		panic(err)
	}

	for {
		a := NewAgent(config, driver)
		err = a.WaitForMatch()
		if err != nil {
			panic(err)
		}
	}
}

type Agent struct {
	name, challenge, match string
	boards                 int
	slow                   bool
	driver                 tictactoe.Driver
}

func NewAgent(cfg Config, d tictactoe.Driver) Agent {
	return Agent{
		name:   cfg.Bot,
		boards: cfg.Boards,
		slow:   cfg.RunSlowly,
		driver: d,
	}
}

func (a Agent) WaitForMatch() error {
	a.challenge = xid.New().String()
	state, err := a.driver.ChallengeState(a.name, a.challenge)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			state, err = a.driver.ChallengeState(a.name, a.challenge)
			if err != nil {
				return err
			}

			if state.Active {
				if state.Match == "" {
					// waiting for a match
				} else if state.Confirmed {
					fmt.Printf("%s Confirmed!\n", time.Now().Format(logTimeFormat))

					// we're done
					ticker.Stop()
					return nil
				} else {
					fmt.Printf("%s Match found\n", time.Now().Format(logTimeFormat))

					// match found, confirm
					a.match = state.Match
					err = a.driver.Confirm(a.name, a.challenge, a.match)
					if err == tictactoe.ErrChallengeNotFound {
						state.Active = false
					} else if err != nil {
						// oops confirmed at a bad time
						return err
					}
				}
			} else {
				//fmt.Printf("%s Creating new challenge\n", time.Now().Format(logTimeFormat))
				err := a.driver.Challenge(a.name, a.boards, "", a.challenge)
				if err != nil {
					return err
				}
			}
		}
	}
}
