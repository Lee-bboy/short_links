package common

import (
	"shortlinks/config"
	"strconv"
	"sync"

	"github.com/garyburd/redigo/redis"
)

var RedisPool *redis.Pool

//确保每个自增键只能对应一个AutoIncrementId对象
var autoIncrementIds map[string]*AutoIncrementId
var autoMutex *sync.Mutex

func init() {
	var err error
	host := config.GetConf("redis.host", "")
	port := config.GetConf("redis.port", "")
	pass := config.GetConf("redis.pass", "")
	dbString := config.GetConf("redis.db", "")
	poolSize := config.GetConf("redis.pool", "")
	db, err := strconv.Atoi(dbString)
	if err != nil {
		panic(err)
	}
	pool, err := strconv.Atoi(poolSize)
	if err != nil {
		panic(err)
	}

	RedisPool = redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", host+":"+port, redis.DialPassword(pass), redis.DialDatabase(db))
	}, pool)

	autoIncrementIds = make(map[string]*AutoIncrementId)
	autoMutex = new(sync.Mutex)
}

type AutoIncrementId struct {
	key   string //redis里面对应的key值
	mutex *sync.Mutex
}

func NewAutoIncrementId(key string) *AutoIncrementId {
	autoMutex.Lock()
	defer autoMutex.Unlock()

	_, ok := autoIncrementIds[key]
	if !ok {
		autoIncrementIds[key] = newAutoIncrementId(key)
	}

	return autoIncrementIds[key]
}

func newAutoIncrementId(key string) *AutoIncrementId {
	autoIncrementId := new(AutoIncrementId)
	autoIncrementId.key = key
	autoIncrementId.mutex = new(sync.Mutex)

	conn := RedisPool.Get()
	conn.Do("SETNX", key, 1)
	conn.Close()

	return autoIncrementId
}

func (this *AutoIncrementId) Get() uint64 {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	conn := RedisPool.Get()
	reply, _ := conn.Do("GET", this.key)
	conn.Do("INCR", this.key)
	conn.Close()

	idBytes, _ := reply.([]byte)
	id, _ := strconv.ParseUint(string(idBytes), 10, 64)
	return id
}
