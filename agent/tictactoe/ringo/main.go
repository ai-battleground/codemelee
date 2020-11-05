package main

import (
	"fmt"
	tictactoe "github.com/ai-battleground/codemelee/client/tictactoe/redis"
	"github.com/rs/xid"
	"sync"
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
		go a.Play()
	}
}

type Agent struct {
	name, challenge, match string
	boards                 int
	slow                   bool
	driver                 tictactoe.Driver
	token                  string
}

func NewAgent(cfg Config, d tictactoe.Driver) Agent {
	return Agent{
		name:   cfg.Bot,
		boards: cfg.Boards,
		slow:   cfg.RunSlowly,
		driver: d,
		token:  xid.New().String(),
	}
}

func (a *Agent) WaitForMatch() error {
	a.challenge = xid.New().String()
	state, err := a.driver.ChallengeState(a.name, a.challenge)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(500 * time.Millisecond)
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
					fmt.Printf("%s Confirmed %s!\n", time.Now().Format(logTimeFormat), state.Match)

					// we're done
					ticker.Stop()
					return nil
				} else {
					fmt.Printf("%s Match found\n", time.Now().Format(logTimeFormat))

					// match found, confirm
					a.match = state.Match
					err = a.driver.Confirm(a.name, a.challenge, a.match, a.token)
					if err == tictactoe.ErrChallengeNotFound {
						state.Active = false
					} else if err != nil {
						// oops confirmed at a bad time
						return err
					}
				}
			} else {
				//fmt.Printf("%s Creating new challenge\n", time.Now().Format(logTimeFormat))
				err := a.driver.Challenge(a.name, a.boards, "", a.challenge, a.token)
				if err != nil {
					return err
				}
			}
		}
	}
}

// play until game over
func (a *Agent) Play() {
	fmt.Printf("%s Starting play %s.\n", time.Now().Format(logTimeFormat), a.match)

	ticker := time.NewTicker(50 * time.Millisecond)
	var busy bool
	var lock sync.Mutex
	for {
		select {
		case <-ticker.C:
			lock.Lock()
			if busy {
				// skip this cycle
				lock.Unlock()
			} else {
				busy = true
				lock.Unlock()
				done := a.actOnce()
				if done {
					ticker.Stop()
					return
				}
				lock.Lock()
				busy = false
				lock.Unlock()
			}

		}
	}
}

func (a *Agent) actOnce() (done bool) {
	o := a.driver.Observe(a.name, a.challenge)
	if o.Error != nil {
		fmt.Printf("%s [%s] Observe error %v\n", time.Now().Format(logTimeFormat), a.match, o.Error)
		return false
	}
	actions := a.Act(o)
	fmt.Printf("%s [%s] Action:%s\n", time.Now().Format(logTimeFormat), a.match, actions)

	if actions != "" {
		err := a.driver.Act(a.name, a.challenge, a.token, actions)
		if err != nil {
			fmt.Printf("%s [%s] Error acting %v\n", time.Now().Format(logTimeFormat), a.match, err)
			return true
		}
		if o.State == "Done" {
			// game over
			return true
		}
	}
	return false
}
