# MyEnvelope

#### 介绍
红包雨项目

#### 软件架构

该项目使用go mod自动进行包管理

**目录说明：**

- model包中定义了保存在数据库中的红包实体
- algo包中实现了随机生成红包的函数，一些可配置的数据也定义在其中
- api包中包含处理不同请求的具体实现
- dao包中实现了对数据库的CRUD操作，以及建表语句

**已实现功能：**

- 一些参数可自定义（暂无持久化代码）；
- get_wallet_list可显示红包总额、根据是否拆开红包来判断是否显示金额、红包按时间排序
- 红包ID使用UUID（36个字符），和用户ID一起，都选用字符串作为id（不知道是否会影响效率）
- 等等......


#### 安装教程
测试环境：

go version go1.17.2 windows/amd64

mysql 5.7

**换源**

> go env -w GO111MODULE=on
> go env -w GOPROXY=https://goproxy.cn,direct

#### 使用说明

以下提交方式均为post：

1.  open接口(.../open)，提供uid，envelope_id 
2.  snatch接口(.../snatch)，提供uid
3.  get_wallet_list接口(.../get_wallet_list），提供uid

可以直接通过使用不同uid运行snatch接口来在数据库中插入数据。

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


#### 特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  Gitee 官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解 Gitee 上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是 Gitee 最有价值开源项目，是综合评定出的优秀开源项目
5.  Gitee 官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  Gitee 封面人物是一档用来展示 Gitee 会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)


#### 配置管理

思路：将公共配置写入`configs/default.yaml`，通过传入参数/环境变量方式导入特定环境下的配置信息。在部署流水线上可以通过Dockfile设置环境变量，或者设置`go run`语句实现配置导入。