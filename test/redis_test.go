package test

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func TestRedis(t *testing.T) {
	err := rdb.Set(ctx, "key", "good", time.Second*20).Err()
	if err != nil {
		t.Fatal(err)
	}

	result, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("value", result)

}
