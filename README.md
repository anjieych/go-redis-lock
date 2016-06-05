# go-redis-lock
go redis lock base on github.com/garyburd/redigo

# Usage:
import "github.com/anjieych/go-redis-lock"

lock := redislock.NewRedislock("192.168.200.88:6379", 0, "")
ok, err := lock.Trylock("Key1", "Key1_value1", 1000)
fmt.Println("redislock Trylock: ", ok, err)
err =lock.Unlock("Key1")
fmt.Println("redislock Unlock: ",  err

# Contact Me:
Email:anjieych@126.com
QQ:   272348197
