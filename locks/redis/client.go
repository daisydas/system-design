package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Data struct {
	Client *redis.Client
}

func (d *Data) GetKeyFromRedis(ctx context.Context, server string) error {
	script := `
		local result = redis.call("SETNX", KEYS[1], ARGV[1])
		if result == 1 then
			redis.call("EXPIRE", KEYS[1], ARGV[2])
			return 1
		else
			return 0
		end
`
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			result, err := d.Client.Eval(ctx, script, []string{"occupy"}, server, 50*time.Second).Result()
			if err != nil {
				fmt.Println("redis setnx error:", err)
				continue
			}
			if result.(int64) == 1 {
				return nil
			}

		}
	}
}

func (d *Data) RemoveLock(ctx context.Context, server string) error {

	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
`

	val, err := d.Client.Eval(ctx, script, []string{"occupy"}, server).Result()
	if err != nil {
		fmt.Println("failed to relaese lock")
		return errors.New("failed to relaese lock")
	}

	if val.(int64) == 0 {
		return fmt.Errorf("lock not held by this server")
	}
	return nil
}
