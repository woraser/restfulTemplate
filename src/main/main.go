package main

import (
	"github.com/unrolled/render"
	"log"
	"os"
)

const local string = "LOCAL"

func main() {
	var (
		// environment variables
		env      = os.Getenv("ENV")      // LOCAL, DEV, STG, PRD
		port     = os.Getenv("PORT")     // server traffic on this port
		version  = os.Getenv("VERSION")  // path to VERSION file
	)
	if env == "" || env == local {
		// running from localhost, so set some default values
		env = local
		port = "3001"
		version = "VERSION"
	}
	//初始化配置文件
	//默认文件名称:config.ini
	configParse,cfg_err :=InitConfigParse("config.ini")
	if cfg_err != nil{
		log.Fatal(cfg_err)
	}
	//初始化mysql
	db,db_err:= InitMySqlDataBase(configParse)
	if db_err != nil{
		log.Fatal(db_err)
	}
	//初始化redis
	redisClient,redis_err :=InitRedisClient(configParse)
	if redis_err != nil{
		log.Fatal(redis_err)
	}
	// initialse application context
	// 初始化应用上下文
	ctx := AppContext{
		Render:  render.New(),
		Version: version,
		Env:     env,
		Port:    port,
		DB:      db,
		RedisClient:redisClient,
		ConfigParse:configParse,
	}
	// start application
	StartServer(ctx)
}
