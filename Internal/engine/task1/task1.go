package task1

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Task1 struct {
	Ctx        context.Context
	TestString string
	Rdb        redis.ClusterClient
}

func New(ctx context.Context, rdb redis.ClusterClient) (*Task1, error) {
	task1 := &Task1{
		Ctx:        ctx,
		TestString: "hello",
		Rdb:        rdb,
	}
	return task1, nil
}
func (task1 *Task1) Start() {
	fmt.Println(task1.TestString)
}
