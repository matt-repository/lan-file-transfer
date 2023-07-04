# lan-file-transfer
局域网传输文件

### 编译步骤：

1、将前端的dist文件夹拷贝到后台项目目录中;

2、执行命令,将前端文件生成二进制文件;

```shell
go-bindata -o asset/asset.go dist/...
```

3、将asset里的包名改成asset（生成默认的时候是main）

4、编译后台go项目

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lan-file-transfer_linux
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o lan-file-transfer_windows.exe
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o lan-file-transfer_mac
```

