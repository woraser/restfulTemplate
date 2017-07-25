package main

import (
	"github.com/unrolled/render"
	"os"
	"log"
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
	// reading version from file
	version = "RELEASE"

	db,db_err:= InitMySqlDataBase("root:root@tcp(127.0.0.1:3306)/awesome?charset=utf8")
	if db_err != nil{
		log.Fatal(db_err)
	}
	redisClient,redis_err :=InitRedisClient("10.2.1.239:6379",0,200,2000)
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
	}
	// start application
	StartServer(ctx)
}
