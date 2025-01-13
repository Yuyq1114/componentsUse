package task1

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Task1 struct {
	Ctx        context.Context
	TestString string
	Rdb        *redis.ClusterClient
	Pg         *gorm.DB
	DorisDb    *gorm.DB
}

func New(ctx context.Context, rdb *redis.ClusterClient, pg *gorm.DB, dorisDb *gorm.DB) (*Task1, error) {
	task1 := &Task1{
		Ctx:        ctx,
		TestString: "hello",
		Rdb:        rdb,
		Pg:         pg,
		DorisDb:    dorisDb,
	}
	return task1, nil
}
func (task1 *Task1) Start() {
	fmt.Println(task1.TestString)
}
