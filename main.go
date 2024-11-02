package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func NoCache() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// ctx.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		// ctx.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		ctx.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		ctx.Next()
	}
}

func main() {
	serverHost := flag.String("server", "http://0.0.0.0:9999", "server url")
	redisHost := flag.String("redis", "localhost:6379", "redis host")
	db1 := flag.Int("DB1", 5, "db name")
	db2 := flag.Int("DB2", 6, "db name")

	// 解析命令行参数
	flag.Parse()

	rdb1 := redis.NewClient(&redis.Options{
		Addr:     *redisHost, // Redis服务器地址和端口
		Password: "",         // Redis访问密码，如果没有可以为空字符串
		DB:       *db1,       // 使用的Redis数据库编号，默认为0
	})

	rdb2 := redis.NewClient(&redis.Options{
		Addr:     *redisHost, // Redis服务器地址和端口
		Password: "",         // Redis访问密码，如果没有可以为空字符串
		DB:       *db2,       // 使用的Redis数据库编号，默认为0
	})

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(NoCache())

	v1 := router.Group("api")
	{
		v1.GET("/airdrop/evm_total_airdrop", func(c *gin.Context) {
			addressStr := c.Query("address")
			address := common.HexToAddress(addressStr)
			key1 := fmt.Sprintf("address-{%s}", address.String())
			key2 := fmt.Sprintf("bind-{%s}", address.String())
			rdb1.Get(context.Background(), key1)
			rdb2.Get(context.Background(), key2)

			c.JSON(http.StatusOK, gin.H{
				"message": "ok"})
		})
	}

	router.Run(*serverHost)
}
