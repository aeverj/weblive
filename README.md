# webinfo #

**webinfo** 是一个高并发网站信息获取工具，可用于

* 获取到目标相关子域名大量资产后(**支持包括IP段,域名**)，进行存活扫描
* 可自定义header、请求方法、请求体、请求路径、端口，可设定是否跟踪301跳转
* 获取域名解析的ip，识别cdn，轻量级识别指纹、获取标题
* 可以自定义app.json文件，进行自定义的指纹识别，[app.json配置教程](https://github.com/AliasIO/wappalyzer)，[最新版app.json下载](https://github.com/AliasIO/wappalyzer/blob/master/src/technologies.json)

更新记录
* **20230814**：由cdn指纹代替cdn的ip库，单文件编译更便捷

webinfo使用go语言编写

* 发挥`golang`协程优势，快速扫描获取网站全面信息
* 多平台通用

------

### 安装 ###

	git clone https://github.com/aeverj/weblive.git
	cd weblive
	go build weblive.go

### 开始使用

<details>
<summary> 👉 weblive 帮助 👈</summary>

```
Usage of webinfo.exe:
  -H value
        Custom Header
  -M string
        Request Method (default "GET")
  -dataFile string
        The Post data file path
  -follow_redirects
        Follow Redirects
  -iF string
        Load urls from file (default "input.txt")
  -output string
        Output file
  -path string
        Request Path (default "/")
  -ports string
        Custom ports
  -threads int
        Number of threads (default 50)
  -timeout int
        Timeout in seconds (default 3)
```
</details>

#### 直接使用
```
将待扫描目标放到当前目录下input.txt文件中，执行程序
weblive -iF input.txt
```
#### 自定义header
```
weblive -H "X-Forwarded-For:127.0.0.1" -H "X-Originating-IP:127.0.0.1"
```
#### 自定义请求方法 GET|POST ,可支持自定义post数据
```
weblive -M POST -dataFile post数据文件路径
```
#### 自定义请求端口
```
weblive -ports 80,443,8000
```
#### 自定义请求路径
```
weblive -path /admin/login.html
```

结果会放到result文件夹中，网站信息保存为CSV表格

###  扫描结果

| URL                      | Redirect                 | Title                                                       | Status_Code | IP                        | CDN   | Finger                                                       |
| ------------------------ | ------------------------ | ----------------------------------------------------------- | ----------- | ------------------------- | ----- | ------------------------------------------------------------ |
| https://www.baidu.com    | https://www.baidu.com    | 百度一下，你就知道                                          | 200         | 182.61.200.7,182.61.200.6 | false | jQuery                                                       |
| https://github.com       | https://github.com       | The  world’s leading software development platform · GitHub | 200         | 13.250.177.223            | false | Ruby on  Rails,GitHub Pages,Bootstrap                        |
| https://studygolang.com/ | https://studygolang.com/ | 首页 -  Go语言中文网 - Golang中文社区                       | 200         | 59.110.219.94             | false | jQuery,Bootstrap,Google  AdSense,Marked,Gravatar,Nginx,Font Awesome |

### TODO
- [ ] 对存活的网站进行截图
- [ ] 导出结果增加html格式