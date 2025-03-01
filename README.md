# 短链接项目
> 本短链接项目旨在将长 URL 转换为简洁的短链接，方便用户分享和使用。通过该项目，用户可以提交长 URL，系统会生成对应的唯一短链接，访问短链接时会自动重定向到原始的长 URL。

## 功能特性
* 短链接生成：用户可以输入长 URL，系统会生成对应的短链接。
* 唯一标识：每个长 URL 会生成唯一的短链接，确保不会重复。
* 重定向功能：访问短链接时，系统会自动重定向到原始的长 URL。
* 认证与授权：部分功能可能需要用户进行认证和授权，确保数据安全。



## 启动步骤
> 修改各个微服务的配置文件

* docker-compose up
* go run gateway.go
* go run auth.go
* go run convert.go
* go run redict.go

## 技术栈
Go-Zero
MySQL
Redis
jwt
...

## 目录结构
.
├── docker-compose.yaml          # Docker 容器编排配置文件
├── README.md                    # 项目说明文档
├── script                       # 脚本文件目录
│   └── init.sql                 # 数据库初始化 SQL 脚本
└── short_link                   # 短链接项目核心代码目录
    ├── common                   # 公共模块目录
    │   └── etcd                 # Etcd 相关功能模块
    │       ├── delivery_address.go  # Etcd 地址相关代码
    │       └── init.go          # Etcd 初始化代码
    ├── go.mod                   # Go 模块依赖管理文件
    ├── go.sum                   # Go 模块依赖校验文件
    ├── main.go                  # 项目入口文件
    ├── pkg                      # 工具包目录
    │   ├── base62               # Base62 编码相关工具
    │   │   └── base62.go
    │   ├── bloom                # 布隆过滤器相关工具
    │   │   └── bloom.go
    │   ├── connect              # 连接相关工具
    │   │   └── connect.go
    │   ├── enter.go             # 可能是入口相关的工具文件
    │   ├── ip                   # IP 相关工具
    │   │   └── ip.go
    │   ├── jwt                  # JWT 认证相关工具
    │   │   └── jwt.go
    │   ├── mds                  # MD5 相关工具
    │   │   └── md5.go
    │   └── urltool              # URL 处理相关工具
    │       └── urltool.go
    ├── shorturlmapmodel         # 短链接映射模型目录
    │   ├── shorturlmapmodel_gen.go  # 自动生成的短链接映射模型代码
    │   ├── shorturlmapmodel.go  # 短链接映射模型代码
    │   └── vars.go              # 短链接映射模型相关变量文件
    ├── sl_auth                  # 认证服务模块
    │   ├── auth_api             # 认证服务 API 相关
    │   │   ├── auth_api.api     # 认证服务 API 定义文件
    │   │   ├── auth.go          # 认证服务主代码
    │   │   ├── etc              # 认证服务配置文件目录
    │   │   │   └── auth.yaml    # 认证服务配置文件
    │   │   └── internal         # 认证服务内部代码目录
    │   │       ├── config       # 认证服务配置相关代码
    │   │       │   └── config.go
    │   │       ├── handler      # 认证服务请求处理相关代码
    │   │       │   ├── authenticationhandler.go  # 认证处理代码
    │   │       │   ├── loginhandler.go      # 登录处理代码
    │   │       │   ├── routes.go        # 认证服务路由代码
    │   │       │   └── signuphandler.go # 注册处理代码
    │   │       ├── logic        # 认证服务业务逻辑代码
    │   │       │   ├── authenticationlogic.go  # 认证逻辑代码
    │   │       │   ├── loginlogic.go      # 登录逻辑代码
    │   │       │   └── signuplogic.go     # 注册逻辑代码
    │   │       ├── svc          # 认证服务上下文代码
    │   │       │   └── servicecontext.go
    │   │       └── types        # 认证服务数据类型定义代码
    │   │           └── types.go
    │   └── auth_models          # 认证服务模型目录
    │       ├── usermodel_gen.go # 自动生成的用户模型代码
    │       ├── usermodel.go     # 用户模型代码
    │       └── vars.go          # 用户模型相关变量文件
    ├── sl_convert               # 短链接转换服务模块
    │   ├── convert_api          # 短链接转换服务 API 相关
    │   │   ├── convert_api.api  # 短链接转换服务 API 定义文件
    │   │   ├── convert.go       # 短链接转换服务主代码
    │   │   ├── etc              # 短链接转换服务配置文件目录
    │   │   │   └── convert.yaml # 短链接转换服务配置文件
    │   │   ├── internal         # 短链接转换服务内部代码目录
    │   │       ├── config       # 短链接转换服务配置相关代码
    │   │       │   └── config.go
    │   │       ├── handler      # 短链接转换服务请求处理相关代码
    │   │       │   ├── converthandler.go  # 转换处理代码
    │   │       │   └── routes.go    # 短链接转换服务路由代码
    │   │       ├── logic        # 短链接转换服务业务逻辑代码
    │   │       │   └── convertlogic.go  # 转换逻辑代码
    │   │       ├── svc          # 短链接转换服务上下文代码
    │   │       │   └── servicecontext.go
    │   │       └── types        # 短链接转换服务数据类型定义代码
    │   │           └── types.go
    │   │   └── sequence         # 短链接转换服务序列生成相关
    │   │       ├── mysql.go     # 序列生成 MySQL 相关代码
    │   │       └── sequence.go  # 序列生成代码
    │   └── convert_models       # 短链接转换服务模型目录
    │       ├── sequencemodel_gen.go # 自动生成的序列模型代码
    │       ├── sequencemodel.go # 序列模型代码
    │       └── vars.go          # 序列模型相关变量文件
    ├── sl_gateway               # 网关服务模块
    │   ├── error                # 网关服务错误处理目录
    │   │   └── error_response.go # 网关服务错误响应代码
    │   ├── gateway.go           # 网关服务主代码
    │   └── settings.yaml        # 网关服务配置文件
    ├── sl_redict                # 短链接重定向服务模块
    │   ├── redict_api           # 短链接重定向服务 API 相关
    │   │   ├── constants        # 短链接重定向服务常量目录
    │   │   │   └── constants.go # 短链接重定向服务常量代码
    │   │   ├── etc              # 短链接重定向服务配置文件目录
    │   │   │   └── redict.yaml  # 短链接重定向服务配置文件
    │   │   ├── internal         # 短链接重定向服务内部代码目录
    │   │       ├── config       # 短链接重定向服务配置相关代码
    │   │       │   └── config.go
    │   │       ├── handler      # 短链接重定向服务请求处理相关代码
    │   │       │   ├── routes.go    # 短链接重定向服务路由代码
    │   │       │   └── showhandler.go # 重定向展示处理代码
    │   │       ├── logic        # 短链接重定向服务业务逻辑代码
    │   │       │   └── showlogic.go   # 重定向逻辑代码
    │   │       ├── middleware   # 短链接重定向服务中间件目录
    │   │       │   └── clientipmiddleware.go # 客户端 IP 中间件代码
    │   │       ├── svc          # 短链接重定向服务上下文代码
    │   │       │   └── servicecontext.go
    │   │       └── types        # 短链接重定向服务数据类型定义代码
    │   │           └── types.go
    │   │   ├── redict_api.api   # 短链接重定向服务 API 定义文件
    │   │   └── redict.go        # 短链接重定向服务主代码
    │   └── redict_models        # 短链接重定向服务模型目录
    │       ├── shorturlaccesslogmodel_gen.go # 自动生成的短链接访问日志模型代码
    │       ├── shorturlaccesslogmodel.go # 短链接访问日志模型代码
    │       └── vars.go          # 短链接访问日志模型相关变量文件
    └── swap                     # 交换相关模块
        └── trace                # 跟踪相关目录
            └── trace.go         # 跟踪代码