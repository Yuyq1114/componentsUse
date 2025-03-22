package asset

import (
	"context"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"golang.org/x/sync/errgroup"
	"test_component/Internal/model"
	"time"
)

var (
	pollAssetSize = 10
)

type Asset struct {
	Ctx       context.Context
	Eg        *errgroup.Group
	InputChan chan *model.TLV
}

func New(Ctx context.Context, eg *errgroup.Group, InputChan chan *model.TLV) *Asset {
	return &Asset{
		Ctx:       Ctx,
		Eg:        eg,
		InputChan: InputChan,
	}
}

func (asset *Asset) Start() {
	asset.Eg.Go(func() error {
		return asset.GetMergeAsset()
	})
}

func (asset *Asset) GetMergeAsset() (err error) {
	//创建协程池
	funcPool, err := ants.NewPoolWithFunc(pollAssetSize, func(i interface{}) {
		msg := i.(*model.TLV)
		asset.handlerAssetMerge(msg)
	}, ants.WithPreAlloc(true), ants.WithNonblocking(true), ants.WithExpiryDuration(10*time.Second))
	if err != nil {
		fmt.Println("协程池分配失败:", err)
	}
	defer funcPool.Release()

	for {
		select {
		case <-asset.Ctx.Done():
			fmt.Println("asset  mission end")
			return nil
		case msg := <-asset.InputChan:
			for {
				err := funcPool.Invoke(msg)
				if err != nil {
					fmt.Println("协程池执行失败:", err)
				}
				//后续更新考虑协程池满的情况
			}
		}
	}
	return err
}

func (asset *Asset) handlerAssetMerge(msg *model.TLV) {
	fmt.Println(msg)
	return
}
