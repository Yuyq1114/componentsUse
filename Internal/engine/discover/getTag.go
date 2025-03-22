package discover

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"sync"
	"test_component/Internal/model"
)

type Tag struct {
	Ctx       context.Context
	Eg        *errgroup.Group
	InputChan chan *model.TLV
	//OutputChan chan
	Pg *gorm.DB
}

func New(Ctx context.Context, eg *errgroup.Group, InputChan chan *model.TLV, Pg *gorm.DB) *Tag {
	return &Tag{
		Ctx:       Ctx,
		Eg:        eg,
		InputChan: InputChan,

		Pg: Pg,
	}
}

func (tag *Tag) Start() {
	tag.Eg.Go(func() error {
		return tag.GetTag()
	})
}

func (tag *Tag) GetTag() (err error) {
	// 通过协程池，创建多个协程对获取的请求进行打标
	//打标方式采用正则匹配
	for {
		select {
		case <-tag.Ctx.Done():
			fmt.Println("GetTag Done")
			return nil
		case msg := <-tag.InputChan:
			//匹配url
			err = tag.UrlTag(msg)
			if err != nil {
				fmt.Println("url tag error")
			}
			//匹配header
			err = tag.HeaderTag(msg)
			if err != nil {
				fmt.Println("header tag error")
			}
			//匹配body
			err = tag.BodyTag(msg)
			if err != nil {
				fmt.Println("body tag error")
			}
		}

	}

	return nil
}

func (tag *Tag) UrlTag(msg *model.TLV) (err error) {
	return
}

func (tag *Tag) HeaderTag(msg *model.TLV) (err error) {

	//打标多个地方
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	result := make(map[string]string)
	wg.Add(4)
	processTag := func(position string) {
		defer wg.Done()
		data, _ := MatchDataPcre2(position)
		mu.Lock()
		result[position] = data
		mu.Unlock()
	}

	go processTag("1")
	go processTag("2")
	go processTag("3")
	go processTag("4")

	wg.Done()

	return
}

func MatchDataPcre2(position string) (str string, err error) {
	return "123", nil

}

func (tag *Tag) BodyTag(msg *model.TLV) (err error) {
	return
}
