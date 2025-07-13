# 酒店管理系统

## 项目简介
本项目是一个基于Go语言开发的酒店管理系统，提供酒店房间管理、客人入住登记、预订管理等功能，帮助酒店实现数字化管理。系统采用前后端分离的架构，后端使用Gin框架提供RESTful API服务，前端可以通过这些API进行交互。

## 技术栈

### 后端
- 编程语言：Go 1.24.3
- Web框架：Gin 1.10.1
- ORM框架：GORM
- 数据库：MySQL
- 配置管理：Viper
- 认证：JWT (JSON Web Tokens)
- 其他：CORS支持、文件上传处理

### 开发工具
- GoLand IDE

## 项目结构
```
hotel-management-system/
├── config/              # 配置文件和配置相关代码
│   ├── config.go        # 配置加载逻辑
│   ├── config.yml       # 配置文件
│   └── db.go            # 数据库连接配置
├── controllers/         # 控制器，处理HTTP请求
│   ├── admin.go         # 管理员相关控制器
│   ├── file.go          # 文件上传控制器
│   ├── guest.go         # 客人管理控制器
│   ├── role.go          # 角色管理控制器
│   ├── room.go          # 房间管理控制器
│   └── roomType.go      # 房间类型控制器
├── global/              # 全局变量和常量
│   └── global.go        # 全局变量定义
├── middlewares/         # 中间件
│   └── auth.go          # 认证中间件
├── models/              # 数据模型
│   ├── guest.go         # 客人模型
│   ├── img.go           # 图片模型
│   ├── reside.go        # 入住记录模型
│   ├── roles.go         # 角色模型
│   ├── rooms.go         # 房间模型
│   ├── roomTypes.go     # 房间类型模型
│   └── users.go         # 用户模型
├── routers/             # 路由配置
│   └── router.go        # 路由注册
├── utils/               # 工具函数
│   └── utils.go         # 通用工具函数
├── uploads/             # 上传文件存储目录
├── go.mod               # Go模块依赖
├── go.sum               # Go模块校验和
└── main.go              # 程序入口
```

## 功能特性

### 用户管理
- 用户注册与登录
- 基于角色的权限控制
- JWT认证

### 房间管理
- 房间信息的增删改查
- 房间类型管理
- 房间状态追踪（空闲、已预订、已入住等）

### 客人管理
- 客人信息登记与管理
- 客人入住记录管理
- 客史档案查询

### 预订与入住
- 房间预订
- 入住登记
- 退房处理

### 文件管理
- 图片上传功能（如房间照片）

## 安装与运行

### 前置条件
- Go 1.24.3或更高版本
- MySQL数据库

### 获取代码
```bash
git clone https://github.com/kingleonardo999/HMS-backend
cd hotel-management-system
```

### 配置
编辑`config/config.yml`文件，填入正确的数据库连接信息和应用端口：
```yaml
app:
  port: "8080"  # 应用运行端口
  name: "hotel-management-system"
database:
  host: "localhost"
  port: "3306"
  username: "root"
  password: "your_password"
  dbname: "hotel_db"
```

### 安装依赖
```bash
go mod download
```

### 运行
```bash
go run main.go
```
或构建后运行：
```bash
go build
./hotel-management-system
```

### API访问
服务启动后，API接口将在配置的端口上可用（默认8080）：
- http://localhost:8080/admin/login - 管理员登录
- http://localhost:8080/admin/register - 管理员注册
- 其他API端点...

## API文档

### 认证类API
- `POST /admin/register` - 管理员注册
- `POST /admin/login` - 管理员登录

### 房间管理API
- `GET /rooms` - 获取房间列表
- `GET /rooms/:id` - 获取特定房间详情
- `POST /rooms` - 创建新房间
- `PUT /rooms/:id` - 更新房间信息
- `DELETE /rooms/:id` - 删除房间

### 客人管理API
- `GET /guests` - 获取客人列表
- `GET /guests/:id` - 获取特定客人详情
- `POST /guests` - 添加新客人
- `PUT /guests/:id` - 更新客人信息

### 更多API端点详见控制器和路由文件

## 贡献指南
欢迎提交问题报告和功能请求。如果您想贡献代码，请遵循以下步骤：
1. Fork 仓库
2. 创建特性分支 (`git checkout -b feature/your-feature`)
3. 提交更改 (`git commit -m 'Add some feature'`)
4. 推送到分支 (`git push origin feature/your-feature`)
5. 创建Pull Request

## 许可证
[MIT](./LICENSE)