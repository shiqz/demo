## 开发文档说明
需求明细：  
[需求文档-1.md](./docs/初级需求-1.md)  
[需求文档-2.md](./docs/初级需求-2.md)  
[需求文档-3.md](./docs/初级需求-3.md)  

## 快速部署
```bash
docker-compose up
```

## 命令行工具
```bash
# 创建管理员账号示例
# -e: 邮箱账号，-p: 登录密码，-r: 角色
docker exec dc_api /app admin create -e=admin@test.com -p=Test123.. -r=manager

# 修改角色
docker exec dc_api /app admin update-role -e=admin@test.com -r=test

# 帮助
docker exec dc_api /app --help
```

## 通用响应结构
**响应成功**
```json
{
  "status": 200,
  "message": "success",
  "data": {}
}
```

**错误响应**
```json
{
  "status": 1004,
  "message": "该账号已被注册",
  "data": {}
}
```

## 内部错误码说明
| 状态码    | 说明       |
|:-------|:---------|
| `1000` | 请求错误     |
| `1001` | 服务端内部错误  |
| `1002` | 请求接口不存在  |
| `1003` | 用户名或密码错误 |
| `1004` | 该用户名已被注册 |
| `1005` | 客户端未登录授权 |
| `1006` | 账号被禁用    |
| `1007` | 无操作权限    |