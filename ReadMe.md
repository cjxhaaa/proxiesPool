代理池
======
##配置文件说明：

>server：

* port  监听端口
* username  未使用
* password  未使用

>settings

* timeout  代理有效时间
* score\_interval 每次获取代理减分策略
* pass\_score   代理及格分
* init\_score   代理初始分

>qingting_proxy

* enable 是否开启
* num 每次获取代理数量
* order\_id 加密后的order\_id

##接口说明

###GET:

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