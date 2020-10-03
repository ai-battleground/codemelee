package tetris

import (
	"fmt"
	"github.com/mediocregopher/radix/v3"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Driver struct {
	pool               *radix.Pool
	started            map[string]GameTime
	lastActionResultId radix.StreamEntryID
}

func NewDriver(redisUrl string) (Driver, error) {
	d := Driver{
		started: make(map[string]GameTime),
	}
	u, err := url.Parse(redisUrl)
	if err != nil {
		return d, fmt.Errorf("Invalid redis URL %s: %w", u, err)
	}
	if u.Scheme != "redis" {
		return d, fmt.Errorf("Invalid redis URL: scheme %s not recognized", u.Scheme)
	}
	if len(u.Port()) == 0 {
		u.Host = u.Hostname() + ":6379"
	}
	d.pool, err = radix.NewPool("tcp", u.Host, 3)
	return d, err
}

func (d Driver) Observe(game string) Observation {
	var o map[string]string
	result := Observation{}
	err := d.pool.Do(radix.Cmd(&o, "HGETALL", fmt.Sprintf("observe:tetris:%s", game)))
	//fmt.Printf("-HGETALL %s\n", time.Now().Format(logTimeFormat))
	if err != nil {
		result.Error = err
		fmt.Printf("Redis error %v\n", err)
		return result
	}
	for f, v := range o {
		switch f {
		case "score":
			if score, err := strconv.Atoi(v); err == nil {
				result.Score = score
			}
		case "lines":
			if lines, err := strconv.Atoi(v); err == nil {
				result.Lines = lines
			}
		case "level":
			if level, err := strconv.Atoi(v); err == nil {
				result.Level = level
			}
		case "tick":
			if tick, err := strconv.Atoi(v); err == nil {
				if started, ok := d.started[game]; ok {
					e := time.Now().Sub(started.Time)
					ticksElapsed := tick - started.Ticks
					tickDuration := e.Milliseconds() / int64(ticksElapsed)
					result.Elapsed = time.Duration(tickDuration*int64(tick)) * time.Millisecond
				}
			}
		case "state":
			result.State = v
		case "shelf":
			result.Shelf = []byte(v)
		case "board":
			for _, line := range strings.Split(v, "\n") {
				result.Board = append(result.Board, []byte(line))
			}
		}
	}
	//fmt.Printf("-parsed  %s\n", time.Now().Format(logTimeFormat))
	//d.observeActionResult(game, &result)
	//fmt.Printf("-result  %s\n", time.Now().Format(logTimeFormat))
	return result
}

func (d Driver) Act(game, action string) {
	if action == "Start" {
		fmt.Printf("starting %s\n", game)
		d.startGame(game)
		return
	}
	err := d.appendActionStream(game, action)
	if err != nil {
		fmt.Println(err)
	}
}

func (d Driver) startGame(game string) {
	err := d.appendActionStream(game, "S")
	if err != nil {
		fmt.Printf("Failed to start game %s: %v", game, err)
	}

	d.started[game] = GameTime{
		Ticks: d.Observe(game).Tick,
		Time:  time.Now(),
	}
}

func (d Driver) appendActionStream(game, action string) error {
	return d.pool.Do(
		radix.Cmd(nil, "XADD", fmt.Sprintf("act:tetris:%s", game), "*", "actions", action))
}

func (d Driver) observeActionResult(game string, o *Observation) {
	reader := radix.NewStreamReader(d.pool, radix.StreamReaderOpts{
		Streams: map[string]*radix.StreamEntryID{
			fmt.Sprintf("result:tetris:%s", game): &d.lastActionResultId,
		},
		Count: 1,
	})
	if _, entries, ok := reader.Next(); ok && len(entries) > 0 {
		for _, e := range entries {
			result, ok := e.Fields["result"]
			if ok && result != "OK" {
				o.Error = fmt.Errorf("action caused error %v", result)
			}
		}
	}
}
