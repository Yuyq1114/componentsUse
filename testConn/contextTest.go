package testConn

import (
	"context"
	"fmt"
	"time"
)

func ContTest(ctx context.Context) {
	fmt.Println(`hello`)
	time.Sleep(20 * time.Second)

}
