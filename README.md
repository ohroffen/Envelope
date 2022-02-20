# Envelope

# 1. 项目介绍
本项目实现了一个支持高并发的抢红包后端系统，负责处理前端的抢红包、拆红包、获取红包列表等业务请求。

本仓库为处理请求的业务逻辑模块。

## 其他模块
[配置模块](https://github.com/cmhzc/envelope-manager)

[持久化模块](https://github.com/cmhzc/envelope-writer)

## 目录说明：

`entity`文件夹中定义了保存在数据库中的红包实体。

`api`文件夹中包含处理不同业务请求的具体实现。

`redis`文件夹中包含`redis`初始化相关代码。

`mq`文件夹中包含发送红包消息到`Kafka`的相关代码。

# 2. 接口说明

该项目`URL`访问请求方式均为`POST`，因此请求参数均包含在请求体中。

## 抢红包

**请求`URL`：**`.../snatch`

**请求参数：**

```json
{
	"uid": 123
}
```

**响应数据：**

```json
{
	"code": 0,
	"msg":  "success", 
	"data": {
        	"envelope_id": 123,
        	"max_count":   5,
        	"cur_count":   3,
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
    "uid": 123,
    "envelope_id": 123
}
```

**响应数据：**

```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "value": 50
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
    "uid": 123
}
```

**响应数据：**

```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "amount": 112,
        "envelope_list": [
            {
                "envelope_id": 123,
                "value": 50,
                "opened": true,
                "snatch_time": 1634551711
            },
            {
                "envelope_id": 123,
                "opened": false,
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
