<template>
    <div  style="width: 90%;left: 5%;position:relative;">
     <Modal
        v-model="openQrCode"
        width="250"
        title="扫描二维码进入">
        <Select v-model="url" style="width:200px; margin: 0.675rem;" @on-change="updateQrCode">
            <Option v-for="item in urls" :value="item" :key="item">{{ item }}</Option>
        </Select>
        <div style="text-align:center;left: 8px;position:relative;" id="qrcode"></div> 
        <div slot="footer"></div>
    </Modal>
   
    <br/>
    <Button  v-if="!isMobile" style="width: 15%"  type="info" shape="circle" @click="openQrCodeModal">显示系统二维码</Button>
    <br/>
    <img style="width: 15%" src="../assets/tranfer.png" >
    
    <br>
    <Input v-model="selectKey" size="large" placeholder="请输入你想查询的文件名..."  style="width: 60%; margin: 0.675rem;" @on-enter="select_file"></Input> 
    <Button type="primary" @click="select_file">查询</Button> 
   
   
    <br/>
    <Upload
        :on-success="getData"
        multiple
        type="drag"
        action="/api/uploadFile">
        <div style="padding: 1.25rem ;background-color:coral;">
            <Icon type="ios-cloud-upload" size="100" style="color: #8F8F9E"></Icon>
            <p class="p1">点击或将文件拖拽到这里上传</p>
        </div>
    </Upload>
    <br>
    <Table  stripe border :columns="columns1" :data="data1"></Table>
    <Page     :total="tablePage.total"  :current="tablePage.pageIndex"   :page-size-opts="itemsPerPages" 
        show-sizer  
        show-total  
        style="text-align: right;margin: 0.675rem;"
        @on-change="handlePage" @on-page-size-change='handlePageSize'
        ></Page>
    </div>
   
</template>

<script>

import { Button } from 'iview';
import QRCode from 'qrcodejs2'
export default {
    data () {
    let _this=this;
    return {
            selectKey:"",
            tablePage:{
                total:100,
                pageIndex:1,
                pageSize:10,
              
            },
            itemsPerPages:[5,10,20,30,40],
            urls:[],
            url:"",
            openQrCode:false,
            isMobile:false,
            qrcode:"",
            hostPort:"",//测试的时候需要切换ip
            //hostPort:"http://192.168.1.103:9999",//测试的时候需要切换ip
            columns1: [
                {
                    title: '文件名',
                    align:'center',
                    key: 'filename',
                    render: (h, params) => {
                            return h('span', {
                                        style: {
                                            
                                            margin:'2px',
                                            border:'',
                                            color: '#FF7D41',
                                            cursor: 'pointer'
                                        },
                                        on: {
                                            click: () => {
                                                var value=params.row.filename
                                                _this.download_file(value)
                                            }
                                        }
                            }, params.row.filename);
                    }

                },
                {
                    title: '上传时间',
                    key: 'createtime',
                    align:'center',
                    width:100,
                },
                {
                    title: '操作',
                    key: 'action',
                    align:'center',
                    maxWidth:80,
                    render :function (h, params) {
                                return h('div',[h('i-button',{
                                    props:  {
                                                type: 'error',
                                                size: 'small'
                                            },
                                    on:     {
                                                click: () => {
                                                    var value=params.row.filename
                                                    _this.delete_file(value)
                                                }
                                            }
                                },"删除")])
                            }
                }
            ],
            data1: [
                {
                    filename: '图片ssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss.png',
                    createtime:"1"
                },
                {
                    filename: '视频.mp4',
                    createtime:"2"
                },
                {
                    filename: '1.txt',
                    createtime:"3"
                },
                {
                    filename: '身份证.jpg',
                    createtime:""
                }
            ]
        }
    },
    mounted: function () {
        this.getData()
        this.getUrls()
        this.checkIsMobile()
    },
    methods: {
        download_file(fileName) {
            var a = document.createElement("a"); //创建一个<a></a>标签
            a.href = "/data/"+fileName; // 给a标签的href属性值加上地址，注意，这里是绝对路径，不用加 点.
            a.download = fileName; //
            a.style.display = "none"; // 障眼法藏起来a标签
            document.body.appendChild(a); // 将a标签追加到文档对象中
            a.click(); // 模拟点击了a标签，会触发a标签的href的读取，浏览器就会自动下载了
            a.remove(); // 一次性的，用完就删除a标签
        },
        delete_file (value) {
            let _this=this;
            let param={
                fileName:value
            };
            _this.axios.delete('/api/deleteFile',
                {
                    params: param
                }) .then(function (response) {
                    if(response.status==200){
                        _this.getData()
                    }else{
                        _this.$Modal.error({
                            title: '删除',
                            content: response.msg
                        })
                    }
                })
                .catch(function (error) {
                    console.log(error);
            });
            
        },
        select_file(){
            this._data.tablePage.pageIndex = 1;
            this.getData();
        },
        getData(){
            let _this=this;
            //获取数据
            let param={
                key:this._data.selectKey,
                pageIndex:this._data.tablePage.pageIndex,
                pageSize:this._data.tablePage.pageSize
            }
          
            _this.axios.get(_this.hostPort+'/api/getPageListFile',
                {
                    params: param
                })   
                .then(function (response) {
                    if(response.status==200){
                        let _data=[]
                        response.data.data.forEach (function (item, index) {
                        var model = {
                                filename:item.FileName,
                                createtime:_this.numberToTimeStr(item.CreateTime) 
                            }
                        _data.push(model);
                        })
                        _this._data.data1=_data;
                        console.log(_data)
                        _this._data.tablePage.total=response.data.total
                    }else{
                        _this.$Modal.info({
                            title: '查询数据',
                            content: "获取失败"
                        })
                    }
                })
                .catch(function (error) {
                    console.log(error);
                }
            );
        },
        getUrls(){
            let _this=this;
            _this.axios.get(_this.hostPort+'/api/getLocalUrls')   
                .then(function (response) {
                   _this._data.urls=response.data.urls;
                   _this._data.url=response.data.urls[0];
                   _this.getQrCode()
                })
                .catch(function (error) {
                    console.log(error);
                }
            );
        },
        updateQrCode(){
            this.qrcode.clear();
            this.qrcode.makeCode(this.url);
        },

        getQrCode(){
            this.qrcode = new QRCode('qrcode', {
                 text: this.url,
                //text: "https://192.168.1.102:9999",
                width: 200,
                height: 200,
                colorDark : '#000000',
                colorLight : '#ffffff',
                correctLevel : QRCode.CorrectLevel.H,
                
            });
        },
        checkIsMobile() {
            var flag=navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i)
            return flag
        },
        openQrCodeModal(){
            this._data.openQrCode = !this._data.isMobile
        },
        handlePage(value){
            this._data.tablePage.pageIndex = value;
            this.getData();
        },
        handlePageSize(value){
            this._data.tablePage.pageIndex = 1;
            this._data.tablePage.pageSize = value;
            this.getData();
        },
        numberToTimeStr(fmt) {
            const time = new Date(fmt * 1000);
            const Y = time.getFullYear()
            const M = (time.getMonth() + 1).toString().padStart(2, '0')
            const D = time.getDate().toString().padStart(2, '0')
            const h = time.getHours().toString().padStart(2, '0')
            const m = time.getMinutes().toString().padStart(2, '0')
            const s = time.getSeconds().toString().padStart(2, '0')
            return `${Y}-${M}-${D} ${h}:${m}:${s}`
        }
    }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 0.675rem;
}
a {
  color: #42b983;
}

.p1{
    font-family: "宋体","仿宋",sans-serif;/*若电脑不支持宋体，则用仿宋，若不支持仿宋，则在sans-serif中找*/
    font-weight: bold;
    font-size: 150%;
    font-style: italic;
    color: black;/*字体颜色*/
    opacity: 0.7;/*字体的透明度：1：默认样式，0：全透明*/
}
</style>
