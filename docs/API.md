# API文档

## 规范

### 错误响应
```json
{
    "success": false,
    "message": "错误信息"
}
```

### 成功响应
```json
{
    "success": true,
    "message": "执行增加、删除、更新等操作时的成功信息", 
    "data": "查询时返回的数据",
    "total": "分页查询时的数据总数"
}
```

### 键值规范

统一采用小驼峰命名, 包括url查询参数、post请求体

默认键值非空, 除非特别说明

## 角色

### 添加角色
POST /role/add
```json
{
    "name": "Role 1"
}
```

### 删除角色
POST /role/delete
```json
{
    "roleId": 1
}
```

### 更新角色信息
POST /role/update
```json
{
    "roleId": 1,
    "name": "Role 1"
}
```

### 获取角色列表
GET /role/list

### 获取角色信息
GET /role/getOne?roleId=\<int\>

## 用户

### 用户登录
POST /admin/login
```json
{
    "username": "admin",
    "password": "password"
}
```

### 添加用户
POST /admin/add
```json
{
    "loginId": "user1",
    "loginPwd": "password",
    "name": "User 1",
    "phone": "1234567890",
    "email": "lFk1v@example.com",
    "roleId": 1,
    "imgId": 1
}
```
其中email和imgId可选

### 删除用户
POST /admin/delete
```json
{
    "userId": 1
}
```

### 更新用户信息
POST /admin/update
```json
{
    "loginId": "user1",
    "name": "User 1",
    "phone": "1234567890",
    "roleId": 1,
    "imgId": 1
}
```

### 获取用户列表
GET /admin/list?pageIndex=\<int\>&pageSize=\<int\>&roleId=\<int\>
roleId = 0 时不过滤, 表示获取所有用户

### 获取用户信息
GET /admin/getOne?loginId=\<string\>

### 修改用户密码
POST /admin/resetPwd
```json
{
    "loginId": "user1",
    "loginPwd": "password",
    "newLoginPwd": "newPassword"
}
```

## 房型

### 添加房型
POST /roomType/add
```json
{
    "roomTypeName": "单人间",
    "roomTypePrice": 100,
    "typeDescription": "这是描述",
    "bedNum": 1
}
```

### 删除房型

POST /roomType/delete

```json
{
    "roomTypeId": 1
}
```

### 更新房型信息

POST /roomType/update

```json
{
    "id": 1,
    "roomTypeName": "单人间",
    "roomTypePrice": 100,
    "typeDescription": "这是描述",
    "bedNum": 1
}
```

### 获取房型列表

GET /roomType/list

### 获取房型信息

GET /roomType/detail?roomTypeId=\<int\>

## 房间

### 添加房间

POST /room/add

```json
{
    "roomId": "102",
    "roomTypeId": 1,
    "roomStatusId": 1,
    "roomDescription": "房间描述"
}
```

roomDescription 可以为空

### 删除房间

POST /room/delete

```json
{
    "roomId": 1
}
```

### 更新房间信息

POST /room/update

```json
{
    "roomId": "102",
    "roomTypeId": 1,
    "roomStatusId": 1,
    "roomDescription": "房间描述"
}
```

### 获取房间列表

GET /room/list?pageIndex=\<int\>&pageSize=\<int\>&roomTypeId=\<int\>&roomStatusId=\<int\>

### 获取房间信息

GET /room/detail?roomId=\<int\>

### 获取房间状态列表

GET /room/statusList

## 入住信息

### 添加入住信息

POST /guestRecord/add

```json
{
  "identityId": "123456789012345678",
  "guestName": "张三",
  "guestPhone": "13800138000",
  "roomTypeId": 2,
  "roomId": "101",
  "resideDate": "2024-06-01T14:00:00.000Z",
  "deposit": 500,
  "guestNum": 2
}
```

### 删除入住信息

POST /guestRecord/delete

```json
{
    "id": 1
}
```

### 更新入住信息

POST /guestRecord/update

```json
{
  "id": 1,
  "guestPhone": "13800000000",
  "roomTypeId": 2,
  "roomId": "A101",
  "leaveDate": "2024-06-30T12:00:00.000Z",
  "guestNum": 2
}
```

leaveDate 可空

### 获取入住信息列表

GET /guestRecord/list?pageIndex=\<int\>&pageSize=\<int\>&resideState=\<int\>&guestName=\<string\>

### 获取入住信息

GET /guestRecord/detail?roomTypeId=\<int\>

### 获取入住状态列表

GET /guestRecord/statusList

### 获取房间列表

GET /guestRecord/roomList?roomTypeId=\<int\>

根据房型给出空闲的房间列表

### 结账

POST /guestRecord/checkout

```json
{
    "id": 1,
    "totalMoney": 199
}
```

## 订单

### 添加订单

POST /order/add

```json
{
  "identityId": "1234567890",
  "guestName": "张三",
  "guestPhone": "13800138000",
  "roomTypeId": 1,
  "roomId": "A101",
  "resideDate": "2024-06-01T14:00:00.000Z",
  "leaveDate": "2024-06-05T12:00:00.000Z",
  "guestNum": 2,
  "totalMoney": 2000
}
```

### 更新订单信息

POST /order/update

```json
{
  "orderId": "od1234567890",
  "guestPhone": "13800138000",
  "roomId": "101",
  "leaveDate": "2024-07-01T12:00:00.000Z",
  "guestNum": 2,
  "totalMoney": 800
}
```

### 获取订单列表

GET /order/list?pageIndex=\<int\>&pageSize=\<int\>&guestName=\<string\>

### 获取订单信息

GET /order/detail?id=\<int\>

### 订单转入住信息

POST /order/live

```json
{
  "orderId": "od1718000000",
  "guestPhone": "13800138000",
  "roomId": "101",
  "leaveDate": "2024-07-01T12:00:00.000Z",
  "guestNum": 2,
  "totalMoney": 800
}
```

## 菜单

### 添加菜单

POST /menu/add

```json
{
  "name": "红烧牛肉面",
  "typeId": "1",
  "price": 28,
  "imgId": 3,
  "desc": "经典川味，香辣可口"
}
```

### 删除菜单

POST /menu/delete

```json
 {
     "id": 1
 }
```

### 更新菜单信息

POST /menu/update

```json
{
  "id": 1,
  "name": "红烧牛肉面",
  "typeId": "2",
  "price": 28,
  "imgId": 5,
  "desc": "经典川味，香辣可口"
}
```

### 获取菜单列表

GET /menu/list?pageIndex=\<int\>&pageSize=\<int\>&typeId=\<int\>

### 获取菜单信息

GET /menu/detail?id=\<int\>

### 获取菜品列表

GET /menu/typeList

## 字典

### 增加字典

POST /dict/add:dictType

```json
{
    "name": "豪华套房"
}
```

根据不同的字典类型增加至对应的数据库

### 删除字典

POST /dict/delete:dictType

```json
{
    "id": 1
}
```

### 更新字典

POST /dict/delete:dictType

```json
{
    "id": 1,
    "name": "空闲"
}
```

### 查看字典列表

GET /dict/list

返回字典类别

### 查看特定字典

GET /dict/:dictType

返回该类别所有的键值对

## 账单

### 获取销售额列表(房型)

GET /billing/list

### 获取房间入住率前三

GET /billing/top3

## 信息

### 获取系统信息

GET /message/list

### 添加信息

POST /message/add

```json
{
    "loginId": "admin",
    "title": "标题",
    "content": "内容"
}
```

### 删除信息

POST /message/delete

```json
{
    "id": 1
}
```

## 文件上传和下载

### 上传图片

POST /uploads/img

```yaml
requestBody:
  required: true
  content:
    multipart/form-data:
      schema:
        type: object
        properties:
          file:
            type: string
            format: binary
            description: 上传的图片文件（JPEG/PNG）
```

### 下载

GET /uploads/:filename
