# Chigua Blog Backend

基于 Go + Gin + PostgreSQL 构建的现代化博客后端服务。

## 技术栈

- **语言**: Go 1.22+
- **框架**: Gin Web Framework
- **数据库**: PostgreSQL 15+
- **ORM**: sqlx
- **认证**: JWT
- **日志**: logrus

## 功能特性

### 用户模块
- ✅ 用户注册
- ✅ 用户登录
- ✅ 获取当前用户信息

### 文章模块
- ✅ 创建文章
- ✅ 查询文章列表（分页）
- ✅ 查询单篇文章
- ✅ 更新文章（支持部分更新）
- ✅ 删除文章
- ✅ 发布文章
- ✅ 修改文章状态

### 分类模块
- ✅ 创建分类（名称判重）
- ✅ 查询所有分类
- ✅ 删除分类

### 标签模块
- ✅ 创建标签（名称判重）
- ✅ 查询所有标签
- ✅ 删除标签

### 评论模块
- ✅ 创建评论
- ✅ 查询文章评论列表
- ✅ 删除评论（仅作者）

## 项目结构

```
chigua-backend/
├── cmd/                    # 命令入口
│   └── server/             # 服务启动入口
├── config/                 # 配置管理
├── database/               # 数据库连接
├── internal/               # 内部模块
│   ├── api/                # API 控制器
│   ├── middleware/         # 中间件
│   ├── model/              # 数据模型
│   ├── router/             # 路由配置
│   ├── service/            # 业务逻辑
│   └── sql/                # SQL 常量定义
├── resources/              # 资源文件
├── utils/                  # 工具函数
│   ├── ip2region/          # IP 地理位置解析
│   ├── jwt/                # JWT 工具
│   └── logger/             # 日志工具
├── config.yaml.example     # 配置文件示例
└── Dockerfile              # Docker 构建文件
```

## 快速开始

### 环境要求

- Go 1.22+
- PostgreSQL 15+

### 安装依赖

```bash
go mod tidy
```

### 配置文件

复制并修改配置文件：

```bash
cp config.yaml.example config.yaml
```

编辑 `config.yaml` 配置数据库连接等信息：

```yaml
server:
  port: 8080

database:
  host: localhost
  port: 5432
  dbname: chigua
  username: postgres
  password: password

jwt:
  secret: your-secret-key
  expires_hours: 24
```

### 启动服务

```bash
go run cmd/server/main.go
```

### 构建生产版本

```bash
go build -o chigua-backend cmd/server/main.go
./chigua-backend
```

## API 接口

### 用户接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/user/register | 用户注册 | 否 |
| POST | /api/user/login | 用户登录 | 否 |
| GET | /api/user/me | 获取当前用户 | 是 |

### 文章接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/article | 创建文章 | 是 |
| GET | /api/article | 获取文章列表 | 否 |
| GET | /api/article/:id | 获取单篇文章 | 否 |
| PUT | /api/article/:id | 更新文章 | 是 |
| DELETE | /api/article/:id | 删除文章 | 是 |
| POST | /api/article/:id/publish | 发布文章 | 是 |
| PUT | /api/article/:id/status | 修改文章状态 | 是 |

### 分类接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/categorie | 创建分类 | 是 |
| GET | /api/categorie | 获取所有分类 | 否 |
| DELETE | /api/categorie/:id | 删除分类 | 是 |

### 标签接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/tag | 创建标签 | 是 |
| GET | /api/tag | 获取所有标签 | 否 |
| DELETE | /api/tag/:id | 删除标签 | 是 |

### 评论接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/comment | 创建评论 | 是 |
| GET | /api/comment/article/:id | 获取文章评论 | 否 |
| DELETE | /api/comment/:id | 删除评论 | 是 |

## 文章状态常量

| 值 | 状态 | 说明 |
|----|------|------|
| 0 | 草稿 | 未发布 |
| 1 | 已发布 | 公开可见 |
| 2 | 已下架 | 撤回发布 |
| 3 | 审核中 | 待审核 |

## 错误码

| 代码 | 描述 |
|------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 开发规范

### 代码结构

1. **API 层**: 处理 HTTP 请求和响应
2. **Service 层**: 业务逻辑处理
3. **Model 层**: 数据模型定义
4. **SQL 层**: SQL 语句常量定义

### 错误处理

使用预定义错误变量统一处理：
- `ErrUserExists` - 用户已存在
- `ErrCategoryExists` - 分类已存在
- `ErrTagExists` - 标签已存在
- `ErrArticleNoPermission` - 无权限操作文章

### SQL 管理

所有 SQL 语句定义在 `internal/sql/` 目录下，按模块划分：
- `article_sql.go` - 文章相关 SQL
- `category_sql.go` - 分类相关 SQL
- `tag_sql.go` - 标签相关 SQL
- `user_sql.go` - 用户相关 SQL
- `comment_sql.go` - 评论相关 SQL

## Docker 部署

```bash
# 构建镜像
docker build -t chigua-backend .

# 运行容器
docker run -p 8080:8080 --name chigua-backend chigua-backend
```

## License

MIT License
