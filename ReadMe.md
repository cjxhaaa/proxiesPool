代理池
======


## 项目说明：

#### 代理ip特点
* 代理时效性短，约一分钟

* 同一个代理使用频次高会被对方服务器封禁，变得不可用
#### 分析
* 根据代理的相关特点，代理池主要应该实现对有效ip的权重分配，代理的过期检测
#### 如何实现？
考虑过使用redis的有序集合作为存储，可以实现有效ip的权重分配，但是redis只能对整个有序集合设置过期时间，无法对每个集合元素做检查。

为了满足代理池的要求，我自己用代码把跳跃表实现了一遍，以满足需求。

即给新获取的代理ip设置一个初始分，每次被使用扣分。集合内元素顺序按分值由高到低排序，每次取排名第一的元素以获得最健康的代理。代理过期则自动删除。
#### 代理池工作流程
1.读取配置文件，将实现的代理爬虫注册到代理获取器，获取器循环运行注册的爬虫不断获取代理

2.通过一条channel将获取的代理同步到代理集合对象中

3.循环检测代理集合元素超时情况，超时则删除

4.使用httpserver提供代理获取api


## 目录说明：

```$xslt
.
├── ReadMe.md
├── cmd
│   ├── pool
│   │   ├── config.ini                // 配置文件
│   │   └── main.go                   // pool入口
│   └── secret
│       └── main.go
├── go.mod
├── go.sum
└── pkg
    ├── config
    │   └── config.go                 // 解析config.ini模块
    ├── crawler
    │   ├── crawler.go                // Crawler接口和parser解析器注册方法
    │   ├── ip3366                     
    │   │   └── crawler.go            // ip3366代理Crawler接口实现
    │   ├── proxies
    │   │   └── proxies.go            // 用来存代理的通用结构体
    │   └── qingting
    │       ├── crawler.go            // 蜻蜓代理Crawler接口实现
    │       └── crawler_test.go 
    ├── secret
    │   └── secret.go
    ├── set
    │   ├── set.go                    // 有序集合的一些方法实现，基于自己实现的skiplist
    │   └── set_test.go
    ├── settings
    │   └── settings.go             
    └── tools
        └── tools.go

```



## 配置文件说明：

>server：

* port  监听端口
* username  未使用
* password  未使用

>settings

* timeout  代理有效时间
* score\_interval 每次获取代理减分策略
* pass\_score   代理及格分
* init\_score   代理初始分

>qingting

* enable 是否开启
* num 每次获取代理数量
* order\_id 加密后的order\_id

## api说明

### GET:

```
1. /ping

params:num

example:

**注：请求头的UserAgent需改为fuckPool，否则返回403

request:
    http://0.0.0.0:8080/ping?num=2

response:
    {
    "message":"success",
    "proxies":["60.168.81.66:27860","223.242.246.78:49974"]
    }

message内容说明：
not enough proxies：没有取到足够的代理
success：获取成功
no proxies:代理为空，而且proxies字段为空值
```
