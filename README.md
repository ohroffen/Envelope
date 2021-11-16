# Envelope

# 一、项目介绍
本项目基于`MySQl、Redis、Kafka`实现了一个支持高并发的抢红包后端系统，负责处理前端的抢红包、拆红包、获取红包列表等业务请求。

## 1、系统架构

<img src="C:\Users\hongt\Pictures\流程图.jpg" alt="流程图"  />



## 2、目录说明：

`entity`文件夹中定义了保存在数据库中的红包实体。

`api`文件夹中包含处理不同业务请求的具体实现。

`redis`文件夹中包含`redis`初始化相关代码。

`mq`文件夹中包含发送红包消息到`Kafka`的相关代码。

# 二、接口说明

该项目`URL`访问请求方式均为`POST`，因此请求参数均包含在请求体中。

## 抢红包

**请求`URL`：**`.../snatch`

**请求参数：**

```json
{
	"uid": 123 // 用户id
}
```

**响应数据：**

```json
{
	"code": 0, // 成功则code=0，否则为其他
	"msg":  "success", 
	"data": {
        "envelope_id": 123, // 红包id
        "max_count":   5,   // 最多抢几次
        "cur_count":   3,   // 当前为第几次抢
    }
}
```

**响应状态码：**

| 注释                     | 业务 code | msg                    |
| ------------------------ | --------- | ---------------------- |
| 成功                     | `0`       | `success`              |
| 未抢到红包               | `1`       | `miss`                 |
| 用户已抢到上限数量红包   | `2`       | `snatch count used up` |
| 系统全部红包已经分发完毕 | `3`       | `no more envelope`     |
| 输入格式错误             | `4`       | `invalid input`        |

## 拆红包

**请求`URL`：**`.../open`

**请求参数：**

```json
{
    "uid": 123, // 用户id
    "envelope_id": 123 // 红包id
}
```

**响应数据：**

```json
{
    "code": 0,        // 成功则code=0，否则为其他
    "msg": "success",
    "data": {
        "value": 50   // 红包金额，以“分”为单位
    }
}
```

**响应状态码：**

| 注释             | 业务 code | msg                      |
| ---------------- | --------- | ------------------------ |
| 成功             | `0`       | `success`                |
| 红包已经开启过   | `1`       | `already opened`         |
| 用户名下无此红包 | `2`       | `envelope doesn't exist` |
| 输入格式错误     | `3`       | `invalid input`          |

## 获取红包列表

**请求`URL`：**`.../get_wallet_list`

**请求参数：**

```json
{
    "uid": 123 // 用户id
}
```

**响应数据：**

```json
{
    "code": 0, // 成功则code=0，否则为其他
    "msg": "success",
    "data": {
        "amount": 112, // 钱包总额，“分”为单位
        "envelope_list": [
            {
                "envelope_id": 123,
                "value": 50, // 红包面值
                "opened": true, // 是否已拆开
                "snatch_time": 1634551711 // 红包获取时间，UNIX时间戳
            },
            {
                "envelope_id": 123,
                "opened": false, // 未拆开的红包不显示value
                "snatch_time": 1634551711 
            }
        ]
    }
}
```

**响应状态码：**

| 注释         | 业务 code | msg             |
| ------------ | --------- | --------------- |
| 成功         | `0`       | `success`       |
| 输入格式错误 | `1`       | `invalid input` |

