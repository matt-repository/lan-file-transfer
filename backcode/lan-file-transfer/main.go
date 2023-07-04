package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"lan-file-transfer/apps"
	"lan-file-transfer/config"
	"lan-file-transfer/router"
	"path/filepath"
)

const (
	defaultPort = 9999
	dataDir     = "data"
)

func init() {
	port := 0
	flag.IntVar(&port, "port", defaultPort, "service port")
	// 寻找指定端口附近空闲的端口
	port = apps.FindFreePort(port)
	_config := &config.Config{
		ServerPort: port,
		DataDir:    dataDir,
	}
	config.Init(_config)
	//创建文件夹
	apps.CreateDir(filepath.Join(apps.GetCurrentDirectory(), config.Get().DataDir))
}

func main() {
	flag.Parse()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	router.Router(r)
	// 启动一个协程，打开浏览器
	go func() {
		apps.OpenUrl()
	}()
	r.Run(fmt.Sprintf(":%d", config.Get().ServerPort))
}
