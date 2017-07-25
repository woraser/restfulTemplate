package main

import (
	"database/sql"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/unrolled/render"
	"time"
)

// AppContext holds application configuration data
type AppContext struct {
	Render  *render.Render
	Version string
	Env     string
	Port    string
	DB      *sql.DB
	RedisClient *redis.Pool
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
func InitMySqlDataBase(dataSourceName string)(*sql.DB,error){
	if dataSourceName == ""{
		dataSourceName = "root:root@tcp(127.0.0.1:3306)/awesome?charset=utf8"
	}
	db, _ := sql.Open("mysql", dataSourceName)
	db.SetMaxIdleConns(100)                 //最大空闲连接数
	db.SetMaxOpenConns(200)               //最大活动连接数
	db.SetConnMaxLifetime(300 * time.Second) //连接超时时间
	return db,nil
}

//init redis client
func InitRedisClient(REDIS_HOST string,REDIS_DB,MaxIdle,MaxActive int)(*redis.Pool,error){

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
