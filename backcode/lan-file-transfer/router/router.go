package router

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"lan-file-transfer/apps"
	"lan-file-transfer/asset"
	"lan-file-transfer/common"
	"lan-file-transfer/config"
	"net/http"
	"os"
)

const (
	indexHtml = "index.html"
	dist      = "dist"
	static    = "static"
	css       = "css"
	js        = "js"
	fonts     = "fonts"
	img       = "img"
)

func Router(r *gin.Engine) {
	// 执行：go-bindata -o asset.go dist/...    将dist文件夹集成  asset文件夹(go代码)
	fsCss := assetfs.AssetFS{Asset: asset.Asset, AssetDir: asset.AssetDir, AssetInfo: nil, Prefix: common.FileJoin(dist, static, css), Fallback: indexHtml}
	fsJs := assetfs.AssetFS{Asset: asset.Asset, AssetDir: asset.AssetDir, AssetInfo: nil, Prefix: common.FileJoin(dist, static, js), Fallback: indexHtml}
	fsFonts := assetfs.AssetFS{Asset: asset.Asset, AssetDir: asset.AssetDir, AssetInfo: nil, Prefix: common.FileJoin(dist, static, fonts), Fallback: indexHtml}
	fsImg := assetfs.AssetFS{Asset: asset.Asset, AssetDir: asset.AssetDir, AssetInfo: nil, Prefix: common.FileJoin(dist, static, img), Fallback: indexHtml}

	r.StaticFS(common.FileJoin(string(os.PathSeparator), static, css), &fsCss)
	r.StaticFS(common.FileJoin(string(os.PathSeparator), static, fonts), &fsFonts)
	r.StaticFS(common.FileJoin(string(os.PathSeparator), static, img), &fsImg)
	r.StaticFS(common.FileJoin(string(os.PathSeparator), static, js), &fsJs)

	r.GET("/", func(c *gin.Context) {
		c.Writer.WriteHeader(200)
		indexPath, _ := asset.Asset(common.FileJoin(dist, indexHtml))
		_, _ = c.Writer.Write(indexPath)
		c.Writer.Header().Add("Accept", "text/html")
		c.Writer.Flush()
	})
	//上传文件
	r.POST("/api/uploadFile", apps.UploadFile)
	//获取文件分页列表
	r.GET("/api/getPageListFile", apps.GetPageListFile)
	//删除文件
	r.DELETE("/api/deleteFile", apps.DeleteFile)
	//获取可访问本程序的url集合
	r.GET("/api/getLocalUrls", apps.GetLocalUrls)

	//数据
	r.StaticFS(common.FileJoin(string(os.PathSeparator), config.Get().DataDir), http.Dir(common.FileJoin(apps.GetCurrentDirectory(), config.Get().DataDir)))
}
