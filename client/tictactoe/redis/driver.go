package tictactoe

import (
	"errors"
	"fmt"
	"github.com/mediocregopher/radix/v3"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const logTimeFormat = "2006-01-02 15:04:05.000"
const ChallengeLifetimeSecs = 30

var ErrMatchNotFound = errors.New("match not found")
var ErrChallengeNotFound = errors.New("challenge not found")
var ErrInternalDataStore = errors.New("internal data store error")

type Driver struct {
	pool *radix.Pool
}

func NewDriver(redisUrl string) (Driver, error) {
	d := Driver{}
	u, err := url.Parse(redisUrl)
	if err != nil {
		return d, fmt.Errorf("invalid redis URL %s: %w", u, err)
	}
	if u.Scheme != "redis" {
		return d, fmt.Errorf("invalid redis URL: scheme %s not recognized", u.Scheme)
	}
	if len(u.Port()) == 0 {
		u.Host = u.Hostname() + ":6379"
	}
	d.pool, err = radix.NewPool("tcp", u.Host, 3)
	return d, err
}

func (d Driver) Challenge(bot string, boards int, opponent string, userKey string, token string) error {
	if bot == "" {
		return fmt.Errorf("bot name is required")
	}
	if boards < 1 {
		return fmt.Errorf("boards must be positive")
	}
	// create or update challenge
	key := fmt.Sprintf("challenge:tictactoe:%s", keySuffix(bot, userKey))
	hsetArgs := []string{
		key,
		"boards", strconv.Itoa(boards),
	}
	if opponent != "" {
		hsetArgs = append(hsetArgs, "opponent", opponent)
	}
	err := d.pool.Do(radix.Pipeline(
		radix.Cmd(nil, "HSET", hsetArgs...),
		radix.Cmd(nil, "EXPIRE", key, strconv.Itoa(ChallengeLifetimeSecs)),
	))
	if err != nil {
		return errInternalDataStore(err)
	}

	return nil
}

func (d Driver) ChallengeState(bot string, userKey string) (state ChallengeState, err error) {
	err = d.pool.Do(radix.Cmd(&state.Active, "EXISTS", fmt.Sprintf("challenge:tictactoe:%s", keySuffix(bot, userKey))))
	if err != nil {
		fmt.Printf("%s error getting challenge %s %v", time.Now().Format(logTimeFormat), keySuffix(bot, userKey), err)
		return state, errInternalDataStore(err)
	}
	if state.Active {
		err = d.pool.Do(radix.Cmd(&state.Match, "HGET",
			fmt.Sprintf("opportunity:tictactoe:%s", keySuffix(bot, userKey)),
			"match",
		))
		if err != nil {
			fmt.Printf("%s error getting opportunity %s %v\n", time.Now().Format(logTimeFormat), keySuffix(bot, userKey), err)
			return state, errInternalDataStore(err)
		}
		if state.Match == "" {
			return state, nil
		}
	} else {
		// check if there's a game
		err = d.pool.Do(radix.Cmd(&state.Match, "GET",
			fmt.Sprintf("game:tictactoe:%s", keySuffix(bot, userKey))))
		if err != nil {
			return state, errInternalDataStore(err)
		}
		if state.Match == "" {
			// nope nothing there
			return state, nil
		} else {
			state.Active = true
			state.Confirmed = true
		}
	}

	err = d.pool.Do(radix.Cmd(&state.Started, "EXISTS", fmt.Sprintf("observe:tictactoe:%s", state.Match)))
	if err != nil {
		fmt.Printf("%s error getting observations %s %v\n", time.Now().Format(logTimeFormat), keySuffix(bot, userKey), err)
		return state, errInternalDataStore(err)
	}
	return state, nil
}

func (d Driver) Confirm(bot, userKey, match, token string) error {
	// check if match still available (opportunity->match) read
	var reply []string
	err := d.pool.Do(radix.Cmd(&reply, "HMGET", fmt.Sprintf("opportunity:tictactoe:%s", keySuffix(bot, userKey)), "match", "token"))
	if err != nil {
		return errInternalDataStore(err)
	} else if len(reply) != 2 {
		fmt.Printf("%s [Confirm] error getting opportunity %s", time.Now().Format(logTimeFormat), keySuffix(bot, userKey))
		return errInternalDataStore(fmt.Errorf("wrong number of elements in reply %d", len(reply)))
	}
	if reply[0] == "" {
		fmt.Printf("%s [Confirm] No opportunity found %s\n", time.Now().Format(logTimeFormat), keySuffix(bot, userKey))
		return errMatchNotFound(keySuffix(bot, userKey))
	}
	if match != reply[0] {
		fmt.Printf("%s [Confirm] Match %s no longer valid for %s\n", time.Now().Format(logTimeFormat), match, keySuffix(bot, userKey))
	}
	if token != reply[1] {
		fmt.Printf("%s [Confirm] WARNING token has changed during challenge process (%s)\n", time.Now().Format(logTimeFormat), keySuffix(bot, userKey))
	}

	// confirm challenge (opportunity->token, challenge->match) write
	var valid bool
	err = d.pool.Do(radix.Cmd(&valid, "EXISTS", fmt.Sprintf("challenge:tictactoe:%s", keySuffix(bot, userKey))))
	if err != nil {
		fmt.Printf("%s [Confirm] Error confirming challenge %s\n", time.Now().Format(logTimeFormat), keySuffix(bot, userKey))
		return errInternalDataStore(err)
	}
	if !valid {
		fmt.Printf("%s [Confirm] Challenge missing %s\n", time.Now().Format(logTimeFormat), keySuffix(bot, userKey))
		return ErrChallengeNotFound
	}
	err = d.pool.Do(radix.Cmd(nil, "HSET", fmt.Sprintf("opportunity:tictactoe:%s", keySuffix(bot, userKey)),
		"token", token))
	if err != nil {
		fmt.Printf("%s [Confirm] Error setting token for %s\n", time.Now().Format(logTimeFormat), keySuffix(bot, userKey))
		return ErrInternalDataStore
	}

	err = d.pool.Do(radix.Pipeline(
		radix.Cmd(nil, "HSET",
			fmt.Sprintf("challenge:tictactoe:%s", keySuffix(bot, userKey)),
			"match", match),
		radix.Cmd(nil, "SET",
			fmt.Sprintf("game:tictactoe:%s", keySuffix(bot, userKey)),
			match,
			"EX", "3600"), // expire in an hour
	))
	if err != nil {
		fmt.Printf("%s [Confirm] Error confirming challenge %s\n", time.Now().Format(logTimeFormat), keySuffix(bot, userKey))
		return ErrInternalDataStore
	}
	fmt.Printf("%s [Confirm] Confirmation validated %s\n", time.Now().Format(logTimeFormat), match)
	return nil
}

func (d Driver) Observe(bot, userKey string) Observation {
	match, err := d.getMatch(bot, userKey)
	if err != nil {
		return Observation{Error: err}
	}
	key := fmt.Sprintf("observe:tictactoe:%s", match)
	var exists bool
	fact := make(map[string]string)
	err = d.pool.Do(radix.Cmd(&exists, "EXISTS", key))
	if err != nil {
		return Observation{Error: err}
	}
	if !exists {
		return Observation{Error: fmt.Errorf("no match %s for %s:%s", match, bot, userKey)}
	}
	err = d.pool.Do(radix.Cmd(&fact, "HGETALL", key))
	if err != nil {
		return Observation{Error: err}
	}
	o := Observation{
		State: fact["state"],
		Bot:   bot,
	}
	for k, v := range fact {
		switch {
		case strings.HasPrefix(k, "score:"):
			if score, err := strconv.Atoi(v); err == nil {
				scoreBot := strings.SplitN(k, ":", 2)[1]
				if scoreBot == bot {
					o.Score = score
				} else {
					o.OpponentScore = score
					o.Opponent = scoreBot
				}
			}
		case k == "boards":
			for _, b := range strings.Split(v, "\n") {
				o.Boards = append(o.Boards, []byte(b))
			}
		case k == "round":
			if r, err := strconv.Atoi(v); err == nil {
				o.Round = r
			}
		case k == "turn":
			o.MyTurn = v == bot
		case k == "turnExpires":
			if ms, err := strconv.Atoi(v); err == nil {
				o.MoveExpires = time.Unix(0, int64(ms)*1e6)
			}
		case k == "matchExpires":
			if ms, err := strconv.Atoi(v); err == nil {
				o.MatchExpires = time.Unix(0, int64(ms)*1e6)
			}
		}
	}
	return o
}

func (d Driver) Act(bot, userKey, token, actions string) error {
	match, err := d.getMatch(bot, userKey)
	if err != nil {
		return err
	}
	// act:tictactoe:id (stream)
	err = d.pool.Do(radix.Cmd(nil, "XADD", fmt.Sprintf("act:tictactoe:%s", match), "*",
		"token", token,
		"actions", actions))
	if err != nil {
		fmt.Printf("%s [Act] Error (bot: %s, userKey: %s) %v\n", time.Now().Format(logTimeFormat), bot, userKey, err)
		return errInternalDataStore(err)
	}
	return nil
}

func (d Driver) getMatch(bot, userKey string) (string, error) {
	var match string
	err := d.pool.Do(radix.Cmd(&match, "GET",
		fmt.Sprintf("game:tictactoe:%s", keySuffix(bot, userKey))))
	if err != nil {
		return "", errInternalDataStore(err)
	}
	if match == "" {
		return "", errMatchNotFound(keySuffix(bot, userKey))
	}
	return match, nil
}

func keySuffix(bot, userKey string) string {
	if userKey != "" {
		userKey = ":" + userKey
	}
	return bot + userKey
}

func errChallengeNotFound(key string) error {
	return fmt.Errorf("%w: %s", ErrChallengeNotFound, key)
}

func errMatchNotFound(key string) error {
	return fmt.Errorf("%w: %s", ErrMatchNotFound, key)
}

func errInternalDataStore(cause error) error {
	return fmt.Errorf("%w: %v", ErrInternalDataStore, cause)
}
