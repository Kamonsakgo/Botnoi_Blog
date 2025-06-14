package datasources

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

type RedisConnection struct {
	Context context.Context
	Redis   *redis.Client
}

func NewRedisConnection() *RedisConnection {

	opt, _ := redis.ParseURL(os.Getenv("REDIS_URI"))
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Println("redis error: ", err)
	}

	// ctx := context.Background()

	// Initialize cursor
	// var cursor uint64 = 0
	// for {
	// 	// Scan for keys matching the pattern
	// 	keys, nextCursor, err := rdb.Scan(ctx, cursor, "*", 10).Result()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	// Loop through keys and retrieve their values
	// 	for _, key := range keys {
	// 		value, err := rdb.Get(ctx, key).Result()
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		fmt.Printf("Key: %s, Value: %s\n", key, value)
	// 	}

	// 	// Update cursor for next iteration
	// 	cursor = nextCursor

	// 	// Break the loop if cursor is 0
	// 	if cursor == 0 {
	// 		break
	// 	}
	// }

	return &RedisConnection{
		Context: context.Background(),
		Redis:   rdb,
	}
}

// func main() {
// 	ctx := context.Background()

// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})

// 	err := rdb.Set(ctx, "key", "value", 0).Err()
// 	if err != nil {
// 		panic(err)
// 	}

// 	val, err := rdb.Get(ctx, "key").Result()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("key", val)

// 	val2, err := rdb.Get(ctx, "key2").Result()
// 	if err == redis.Nil {
// 		fmt.Println("key2 does not exist")
// 	} else if err != nil {
// 		panic(err)
// 	} else {
// 		fmt.Println("key2", val2)
// 	}
// 	// Output: key value
// 	// key2 does not exist
// }
