package compo

import (
	"context"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/go-redis/redis/v8"
	"github.com/sony/sonyflake"
	"time"
)

type LockBuilder struct {
	rdb *redis.Client
	sf  *sonyflake.Sonyflake
}

// NewLockBuilder .
func NewLockBuilder(rdb *redis.Client) *LockBuilder {
	return &LockBuilder{rdb: rdb}
}

func (lb *LockBuilder) NewLock(ctx context.Context, name string, seconds int64) LockInterface {
	randId, _ := lb.sf.NextID()
	return &redisLock{
		ctx,
		name,
		convertor.ToString(randId),
		seconds,
		lb.rdb,
	}
}

type LockInterface interface {
	Get() bool
	Block(seconds int64) bool
	Release() bool
	ForceRelease()
}

type redisLock struct {
	context context.Context
	name    string
	owner   string
	seconds int64
	rdb     *redis.Client
}

const releaseLockLuaScript = `
if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end
`

func (l *redisLock) Get() bool {
	return l.rdb.SetNX(l.context, l.name, l.owner, time.Duration(l.seconds)*time.Second).Val()
}

func (l *redisLock) Block(seconds int64) bool {
	timer := time.After(time.Duration(seconds) * time.Second)
	for {
		select {
		case <-timer:
			return false
		default:
			if l.Get() {
				return true
			}
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
}

func (l *redisLock) Release() bool {
	luaScript := redis.NewScript(releaseLockLuaScript)
	result := luaScript.Run(l.context, l.rdb, []string{l.name}, l.owner).Val().(int64)
	return result != 0
}

func (l *redisLock) ForceRelease() {
	l.rdb.Del(l.context, l.name).Val()
}
