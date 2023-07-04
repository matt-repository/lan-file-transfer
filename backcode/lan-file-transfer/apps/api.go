package apps

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/gin-gonic/gin"
	"io/fs"
	"io/ioutil"
	"lan-file-transfer/config"
	"lan-file-transfer/model"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	keyStr       = "key"
	pageIndexStr = "pageIndex"
	pageSizeStr  = "pageSize"
	msg          = "msg"
)

func UploadFile(g *gin.Context) {
	file, err := g.FormFile("file")
	if err != nil {
		g.JSON(http.StatusInternalServerError, map[string]interface{}{
			msg: fmt.Sprintf("FormFile err:%s", err.Error()),
		})
		return
	}
	err = g.SaveUploadedFile(file, filepath.Join(GetCurrentDirectory(), config.Get().DataDir, file.Filename))
	if err != nil {
		g.JSON(http.StatusInternalServerError, map[string]interface{}{
			msg: fmt.Sprintf("SaveUploadedFile err:%s", err.Error()),
		})
		return
	}
	g.JSON(http.StatusOK, map[string]interface{}{
		msg: "upload success",
	})
}

func DeleteFile(g *gin.Context) {
	fileName, ok := g.GetQuery("fileName") //获取查询关键字
	if !ok {
		g.JSON(http.StatusOK, map[string]interface{}{
			msg: "file is not exist！",
		})
		return
	}
	err := os.Remove(filepath.Join(GetCurrentDirectory(), config.Get().DataDir, fileName))
	if err != nil {
		g.JSON(http.StatusInternalServerError, map[string]interface{}{
			msg: fmt.Sprintf("os.Remove err:%s", err.Error()),
		})
		return
	}
	g.JSON(http.StatusOK, map[string]interface{}{
		msg: "delete success",
	})

}

func GetPageListFile(g *gin.Context) {
	//获取参数
	pageIndex, pageSize, key, ok := getParam(g)
	if !ok {
		return
	}
	//获取文件集合
	files, ok := getFiles(g)
	if !ok {
		return
	}
	//关键字过滤
	newFiles := filterFiles(key, files)
	//排序
	sort.Sort(newFiles)
	//分页
	total, data := getPageFiles(pageIndex, pageSize, newFiles)
	g.JSON(http.StatusOK, map[string]interface{}{
		"data":  data,
		"total": total,
	})
}

func GetLocalUrls(g *gin.Context) {
	ips := GetLocalIps()
	urls := make([]string, 0)
	for _, ip := range ips {
		urls = append(urls, fmt.Sprintf("http://%s:%d", ip, config.Get().ServerPort))
	}
	g.JSON(http.StatusOK, map[string]interface{}{
		"urls": urls,
	})
}

//getParam 获取参数 pageIndex、pageSize、key
func getParam(g *gin.Context) (pageIndex, pageSize int, key string, ok bool) {
	pageIndexValue, ok := g.GetQuery(pageIndexStr)
	pageErrStr := fmt.Sprintf("pageIndex,pageSize must be greater 0")
	keyErrStr := fmt.Sprintf("key must be offer")
	if !ok {
		g.JSON(http.StatusBadRequest, map[string]interface{}{
			msg: pageErrStr,
		})
		return 0, 0, "", false
	}
	pageSizeValue, ok := g.GetQuery(pageSizeStr)
	if !ok {
		g.JSON(http.StatusBadRequest, map[string]interface{}{
			msg: pageErrStr,
		})
		return 0, 0, "", false
	}

	keyValue, ok := g.GetQuery(keyStr)
	if !ok {
		g.JSON(http.StatusBadRequest, map[string]interface{}{
			msg: keyErrStr,
		})
		return 0, 0, "", false
	}

	pageIndex, err := strconv.Atoi(pageIndexValue)
	if err != nil || pageIndex <= 0 {
		g.JSON(http.StatusBadRequest, map[string]interface{}{
			msg: pageErrStr,
		})
		return 0, 0, "", false
	}
	pageSize, err = strconv.Atoi(pageSizeValue)
	if err != nil || pageSize <= 0 {
		g.JSON(http.StatusBadRequest, map[string]interface{}{
			msg: pageErrStr,
		})
		return 0, 0, "", false
	}
	return pageIndex, pageSize, keyValue, true
}

//获取文件集合
func getFiles(g *gin.Context) ([]fs.FileInfo, bool) {
	path := filepath.Join(GetCurrentDirectory(), config.Get().DataDir)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		g.JSON(http.StatusInternalServerError, map[string]interface{}{
			msg: fmt.Sprintf("get files %s err:%s", path, err.Error()),
		})
		return nil, false
	}
	return files, true
}

//分页获取文件集合
func getPageFiles(pageIndex, pageSize int, files FilesByModTime) (int, []model.FileModel) {
	data := make([]model.FileModel, 0)
	length := len(files)
	for i := 0; i < pageSize; i++ {
		if length > pageSize*(pageIndex-1)+i {
			name := files[pageSize*(pageIndex-1)+i].Name()
			createTime := files[pageSize*(pageIndex-1)+i].ModTime().Unix()
			data = append(data, model.FileModel{FileName: name, CreateTime: createTime})
		}
	}
	return length, data
}

//关键字过滤文件集合
func filterFiles(key string, files FilesByModTime) FilesByModTime {
	newFiles := make([]os.FileInfo, len(files))
	copy(newFiles, files)
	if key != "" {
		linq.From(files).Where(func(i interface{}) bool {
			file := i.(os.FileInfo)
			return strings.Index(file.Name(), key) >= 0
		}).ToSlice(&newFiles)
	}
	return newFiles
}

type FilesByModTime []os.FileInfo

func (fis FilesByModTime) Len() int {
	return len(fis)
}

func (fis FilesByModTime) Swap(i, j int) {
	fis[i], fis[j] = fis[j], fis[i]
}

func (fis FilesByModTime) Less(i, j int) bool {
	return fis[i].ModTime().After(fis[j].ModTime())
}
