// Copyright 2016 Anjieych. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//Usage:
//	lock := redislock.NewRedislock("192.168.200.88:6379", 0, "")
//	ok, err := lock.Trylock("Key1", "Key1_value1", 1000)
//	fmt.Println("redislock Trylock: ", ok, err)
//  err =lock.Unlock(""Key1", "Key1_value1")
//  fmt.Println("redislock Unlock: ",  err)
package redislock

import (
	"fmt"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Redislock struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	dbNum    int
	password string
}

// Connect to redis.
func (lock *Redislock) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", lock.conninfo)
		if err != nil {
			return nil, err
		}

		if lock.password != "" {
			if _, err := c.Do("AUTH", lock.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", lock.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	lock.p = &redis.Pool{
		MaxIdle:     6,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
	c := lock.p.Get()
	defer c.Close()
}

// New redislock
func NewRedislock(conninfo string, dbNum int, password string) *Redislock {
	lock := &Redislock{conninfo: conninfo, dbNum: dbNum, password: password}
	defer lock.connectInit()
	return lock
}


//Try lock the key with value and set timeout_ms
func (lock *Redislock) Trylock(key, value string, timeout_ms int) (ok bool, err error) {
	ok, err = lock.trylock(key, value, timeout_ms)
	if !ok || err != nil {
		lock = nil
	}

	return
}

// Unlock the key
func (lock *Redislock) Unlock(key string) (err error) {
	c := lock.p.Get()
	defer c.Close()
	_, err = c.Do("DEL", key)
	return
}

// Add TTL to to the key
func (lock *Redislock) AddTimeout(key, value string, px_time_ms int64) (ok bool, err error) {
	c := lock.p.Get()
	defer c.Close()
	ttl_time, err := redis.Int64(c.Do("TTL", key))
	fmt.Println(ttl_time)
	if err != nil {
		log.Fatal("redis get failed:", err)
	}
	if ttl_time > 0 {
		fmt.Println(11)
		_, err := redis.String(c.Do("SET", key, value, "PX", int(ttl_time+px_time_ms)))
		if err == redis.ErrNil {
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

func (lock *Redislock) trylock(key, value string, timeout_ms int) (ok bool, err error) {
	c := lock.p.Get()
	defer c.Close()

	_, err = redis.String(c.Do("SET", key, value, "PX", timeout_ms, "NX"))
	if err == redis.ErrNil {
		// The lock was not successful, it already exists.
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
