# CodeDream 

刷题平台

启动： --> 
## 前置要求

1. **MySQL 数据库**
   - 确保 MySQL 服务已启动
   - 创建数据库：`CREATE DATABASE codedream CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;`

2. **Go 环境**
   - Go 1.24 或更高版本

3. **Node.js 环境**
   - Node.js 18+ (推荐 20+)
   - npm 已安装

## 快速启动

### 方法1: 使用启动脚本（推荐）

1. 编辑 `start.ps1` 文件，修改数据库配置：
   ```powershell
   $env:DB_USER = "root"        # 你的MySQL用户名
   $env:DB_PASS = "yourpassword" # 你的MySQL密码
   $env:DB_HOST = "127.0.0.1"
   $env:DB_PORT = "3306"
   $env:DB_NAME = "codedream"
   ```

2. 运行启动脚本：
   ```powershell
   .\start.ps1
   ```

### 方法2: 手动启动

#### 启动后端

```powershell
# 设置数据库环境变量
$env:DB_USER = "root"
$env:DB_PASS = "yourpassword"
$env:DB_HOST = "127.0.0.1"
$env:DB_PORT = "3306"
$env:DB_NAME = "codedream"

# 启动后端服务
go run cmd/server/main.go
```

#### 启动前端

```powershell
cd fronter
npm run dev
```

## 访问地址

- **前端**: http://localhost:5173
- **后端API**: http://localhost:8080
- **社区页面**: http://localhost:8080/community.html

## 数据库配置说明

项目支持通过环境变量配置数据库连接：

- `DB_USER`: MySQL用户名（默认: root）
- `DB_PASS`: MySQL密码（默认: 空）
- `DB_HOST`: MySQL主机（默认: 127.0.0.1）
- `DB_PORT`: MySQL端口（默认: 3306）
- `DB_NAME`: 数据库名（默认: codedream）

## 常见问题

### 1. 数据库连接失败

**错误**: `Access denied for user 'root'@'localhost'`

**解决方案**:
- 检查MySQL用户名和密码是否正确
- 确认MySQL服务已启动
- 检查数据库是否存在：`CREATE DATABASE codedream;`

### 2. 端口被占用

**错误**: `bind: address already in use`

**解决方案**:
- 检查8080端口是否被占用：`netstat -ano | findstr :8080`
- 关闭占用端口的进程或修改端口号

### 3. 前端依赖安装失败

**解决方案**:
- 删除 `fronter/node_modules` 和 `fronter/package-lock.json`
- 重新运行 `npm install`

## API 测试

启动服务后，可以测试API：

```powershell
# 获取帖子列表
Invoke-WebRequest -Uri "http://localhost:8080/api/community/posts" -Method GET

# 创建帖子
$body = @{
    type = "question"
    content = "这是一个测试帖子"
    tags = @("测试", "问题")
} | ConvertTo-Json
Invoke-WebRequest -Uri "http://localhost:8080/api/community/posts" -Method POST -Body $body -ContentType "application/json"
```

## 功能说明

项目已实现以下功能：

✅ 社区帖子管理（创建、查看、分页）
✅ 评论系统
✅ 投票系统（点赞/点踩）
✅ 收藏功能
✅ 热门帖子
✅ 活跃用户
✅ 每日题目推荐

## 开发说明

- 后端使用 Go + Gin + GORM
- 前端使用 React + Vite
- 数据库使用 MySQL
- 所有模型会自动迁移到数据库

