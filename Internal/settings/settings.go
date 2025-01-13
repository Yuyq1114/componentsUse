package settings

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type RedisConfig struct {
	Addrs          []string      `mapstructure:"address"`
	Password       string        `mapstructure:"password"`
	RouteByLatency bool          `mapstructure:"routeByLatency"`
	DialTimeout    time.Duration `mapstructure:"ddialTimeout"` // 设置连接超时时间
	ReadTimeout    time.Duration `mapstructure:"readTimeout"`  // 设置读超时时间
	WriteTimeout   time.Duration `mapstructure:"writeTimeout"`
}
type PgConfig struct {
	DataSource string `mapstructure:"dataSource"`
}

type DorisConfig struct {
	DataSource string `mapstructure:"dataSource"`
}
type Config struct {
	Redis RedisConfig `mapstructure:"redis"`
	Pg    PgConfig    `mapstructure:"PG"`
	Doris DorisConfig `mapstructure:"Doris"`
}

func InitConfig(configPath string) (config Config, err error) {
	config, err = parsePathFile(configPath)
	//后续可以从环境变量导入
	return config, err
}

func parsePathFile(path string) (cfg Config, err error) {
	//var cfg Config
	viper.SetConfigType("yml")
	viper.SetConfigFile(path)
	if err = viper.ReadInConfig(); err != nil {
		log.Println(err)
	}

	//解析redis
	//addr := []string{"string1", "string2"}
	//password := "hello"
	//routeByLatency := true
	//
	//redis := RedisConfig{
	//	Addrs:          addr,
	//	Password:       password,
	//	RouteByLatency: routeByLatency,
	//}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Println(err)
	}
	//fmt.Println(cfg)
	//解析pg

	//解析

	//最后
	//config.Redis = redis
	return
}
