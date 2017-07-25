package main

import (
	"database/sql"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/unrolled/render"
	"time"
	"github.com/robfig/config"
)

// AppContext holds application configuration data
type AppContext struct {
	Render  *render.Render
	Version string
	Env     string
	Port    string
	DB      *sql.DB
	RedisClient *redis.Pool
	ConfigParse *config.Config
}

// Healthcheck will store information about its name and version
type Healthcheck struct {
	AppName string `json:"appName"`
	Version string `json:"version"`
}

// Status is a custom response object we pass around the system and send back to the customer
// 404: Not found
// 500: Internal Server Error
type Status struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

/**init mysql db mysql_driver:_ "github.com/go-sql-driver/mysql"**/
func InitMySqlDataBase(configParse *config.Config)(*sql.DB,error){
	dataSourceName,_ := configParse.String("mysql","dataSourceName")
	MaxIdle,_ := configParse.Int("mysql","MaxIdle")
	MaxOpen,_ := configParse.Int("mysql","MaxOpen")
	db, _ := sql.Open("mysql", dataSourceName)
	db.SetMaxIdleConns(MaxIdle)                 //最大空闲连接数
	db.SetMaxOpenConns(MaxOpen)               //最大活动连接数
	db.SetConnMaxLifetime((60*60*12) * time.Second) //连接超时时间
	return db,nil
}

//init redis client
func InitRedisClient(configParse *config.Config)(*redis.Pool,error){
	REDIS_HOST,_ := configParse.String("redis","REDIS_HOST")
	REDIS_DB,_ := configParse.Int("redis","REDIS_DB")
	MaxIdle,_ := configParse.Int("redis","MaxIdle")
	MaxActive,_ := configParse.Int("redis","MaxActive")

	RedisClient := &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     MaxIdle,  //最大空闲连接数
		MaxActive:   MaxActive, //最大活跃连接数
		Wait:        true,
		IdleTimeout: 3600 * time.Second, //最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		Dial: func() (redis.Conn, error) { //进行连接
			c, err := redis.Dial("tcp", REDIS_HOST)
			if err != nil {
				return nil, err
			}
			// 选择db
			c.Do("SELECT", REDIS_DB)
			return c, nil
		},
	}
	return RedisClient,nil

}

//初始化配置文件对象 默认:config.ini
func InitConfigParse(fileName string)(*config.Config,error){
	configParse, error := config.ReadDefault(fileName)
	return configParse,error

}
