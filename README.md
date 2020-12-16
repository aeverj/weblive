# webinfo #

**webinfo** 是一个高并发网站信息获取工具，可用于

* 获取到目标相关子域名大量资产后，进行存活扫描
* 获取域名解析的ip，识别cdn，轻量级识别指纹、获取标题

webinfo使用go语言编写

* 发挥`golang`协程优势，快速扫描获取网站信息
* 多平台通用

------

### 安装 ###

	go get github.com/aeverj/webinfo

### 开始使用

 **直接扫描**

```
将需要扫描的域名保存到url.txt文件中，执行
webinfo
```

结果会放到result文件夹中，网站信息保存为Excel表格，不存在cdn的真实ip保存到`ip.txt`文件中

###  扫描结果

| URL                      | Redirect                 | Title                                                       | Status_Code | IP                        | CDN   | Finger                                                       |
| ------------------------ | ------------------------ | ----------------------------------------------------------- | ----------- | ------------------------- | ----- | ------------------------------------------------------------ |
| https://www.baidu.com    | https://www.baidu.com    | 百度一下，你就知道                                          | 200         | 182.61.200.7,182.61.200.6 | false | jQuery                                                       |
| https://github.com       | https://github.com       | The  world’s leading software development platform · GitHub | 200         | 13.250.177.223            | false | Ruby on  Rails,GitHub Pages,Bootstrap                        |
| https://studygolang.com/ | https://studygolang.com/ | 首页 -  Go语言中文网 - Golang中文社区                       | 200         | 59.110.219.94             | false | jQuery,Bootstrap,Google  AdSense,Marked,Gravatar,Nginx,Font Awesome |

