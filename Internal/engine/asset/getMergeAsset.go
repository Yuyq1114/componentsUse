package asset

import (
	"context"
	"fmt"
	"github.com/golang/groupcache/lru"
	"github.com/panjf2000/ants/v2"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"log"
	"net/url"
	"test_component/Internal/model"
	"test_component/Internal/utils"
	"time"
)

var (
	pollAssetSize = 5
	cacheCap      = 10
)

type Asset struct {
	Ctx       context.Context
	Eg        *errgroup.Group
	InputChan chan *model.TLV
	Pg        *gorm.DB
	Redis     *redis.Client
	Doris     *gorm.DB
	Cache     *lru.Cache
}

func New(ctx context.Context, eg *errgroup.Group, inputChan chan *model.TLV, pg *gorm.DB, redis *redis.Client, doris *gorm.DB) *Asset {
	return &Asset{
		Ctx:       ctx,
		Eg:        eg,
		InputChan: inputChan,
		Pg:        pg,
		Redis:     redis,
		Doris:     doris,
		//创建一个缓存，用来存储从pg里查询的host，防止每次都要从数据库读
		//综合考虑采用LFU的淘汰策略，因为某一时间段会固定大量访问某个网站
		//更新策略采用旁路缓存，先更新数据库，在删缓存，不然会数据不一致
		//算了 现在写不出来 使用github.com/golang/groupcache/lru这个吧
		//如果有bigcache的需求也可采用github.com/allegro/bigcache
		Cache: lru.New(cacheCap),
	}
}

// Start 启动asset 当前版本仅实现资产的增删该和同步数据库，目前仅通过host字段进行
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
	}, ants.WithPreAlloc(true), ants.WithExpiryDuration(10*time.Second))
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
			err = funcPool.Invoke(msg)
			if err != nil {
				fmt.Println("协程池执行失败:", err)
			}
			//后续更新考虑协程池满的情况

		}
	}
	return err
}

// handlerAssetMerge 处理资产合并
func (asset *Asset) handlerAssetMerge(msg *model.TLV) {
	//目前采用strings包的东西进行匹配，后续的要做成用户可配置的，要采用正则去匹配
	//提取host字段
	machine, err := utils.NewSnowflake(1)
	parse, err := url.Parse(msg.HttpsData.URL)
	if err != nil {
		fmt.Println("url解析失败:", err)
		return
	}
	//parse.Host

	//
	var mergeData model.MergeHttps
	//如果 已有该host，新增一个，如果没有，创建一个
	if value, ok := asset.Cache.Get(parse.Host); ok {
		////pg中插入一条
		//
		//// 插入数据
		//result := asset.Pg.Create(&mergeData)
		//if result.Error != nil {
		//	fmt.Println("插入数据失败: ", result.Error)
		//} else {
		//	fmt.Println("数据插入成功，MergeId: ", mergeData.MergeId)
		//}
		//将mergeid插入doris的日志表
		result := asset.Doris.Table("https_logs").Where("log_id = ?", msg.LogId).Update("merge_id", value)
		if result.Error != nil {
			log.Fatalf("更新用户失败: %v", result.Error)
		}
		fmt.Printf("成功更新用")
	} else {
		//没查到就查库
		if err = asset.Pg.Table("merge_https").Where("hosts = ?", parse.Host).First(&mergeData).Error; err != nil {
			//没查到
			//先更新数据库在更新缓存
			// 插入pg
			mergeData = model.MergeHttps{
				MergeId: machine.GenerateID(),
				Hosts:   parse.Host, // 示例的 Hosts 数据
				//LogId:      msg.LogId,  // 示例的 LogId
				CreateTime: time.Now(), // 当前时间作为 CreateTime
				UpdateTime: time.Now(), // 当前时间作为 UpdateTime
			}
			result := asset.Pg.Table("merge_https").Create(&mergeData)
			if result.Error != nil {
				fmt.Println("插入数据失败: ", result.Error)
				//插入失败就要返回了 不然也会不一致
				return
			} else {
				fmt.Println("数据插入成功，MergeId: ", mergeData.MergeId)
			}
			//直接更新doris
			//但是doris每5条插入一次，这个地方会由于数据不一致导致插入不了
			result = asset.Doris.Table("https_logs").Where("log_id = ?", msg.LogId).Update("merge_id", mergeData.MergeId)
			if result.Error != nil {
				log.Fatalf("更新用户失败: %v", result.Error)
			}
			fmt.Printf("doris成功更新")
			//更新缓存

			if asset.Cache.Len() > 10 {
				asset.Cache.RemoveOldest()
			}
			asset.Cache.Add(parse.Host, mergeData.MergeId)
		} else {
			//查到了
			//更新缓存
			//直接更新doris
			result := asset.Doris.Table("https_logs").Where("log_id = ?", msg.LogId).Update("merge_id", mergeData.MergeId)
			if result.Error != nil {
				log.Fatalf("更新用户失败: %v", result.Error)
			}
			fmt.Printf("成功更新用")
			//更新pg的updatetime
			result = asset.Pg.Table("merge_https").Where("merge_id", value).Update("update_time", time.Now())
			if result.Error != nil {
				log.Fatalf("更新用户失败: %v", result.Error)
			}
			fmt.Printf("pg成功更新")
			//更新缓存
			if asset.Cache.Len() > 10 {
				asset.Cache.RemoveOldest()
			}
			asset.Cache.Add(parse.Host, mergeData.MergeId)
		}

	}

	return
}

// deleteAssetMerge 删除资产合并
func (asset *Asset) deleteAssetMerge(msg *model.TLV) {
	//删除某一个资产 则删除库里所有关联资产
	//前端点击删除，后端发送redis消息通知，然后引擎负责删除

	//监听redis消息

	//处理删除

	return
}
