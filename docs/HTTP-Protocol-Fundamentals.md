# HTTP 协议基础详解

## 概述

HTTP (HyperText Transfer Protocol) 是万维网的基础协议，用于客户端（如浏览器）和服务器之间的通信。HTTP 协议定义了客户端如何请求资源以及服务器如何响应这些请求。

## 协议版本演进

### HTTP/0.9 (1991)
- 极简设计：只有 GET 方法
- 无头信息、无状态码
- 纯文本响应

**示例：**
```http
GET /index.html
```

### HTTP/1.0 (1996) - RFC 1945
- 引入头信息 (Headers)
- 添加状态码 (Status Codes)
- 支持多种方法：GET, POST, HEAD
- 每个请求建立新连接

### HTTP/1.1 (1997) - RFC 2068/2616/7230
- **持久连接 (Persistent Connections)**: Keep-Alive
- **管道化 (Pipelining)**: 并行请求
- **分块传输编码 (Chunked Transfer Encoding)**
- **缓存机制 (Caching)**
- **主机头 (Host Header)**: 支持虚拟主机

### HTTP/2 (2015) - RFC 7540
- **二进制协议**: 更高效解析
- **多路复用 (Multiplexing)**: 单个连接多个流
- **头部压缩 (Header Compression)**: HPACK
- **服务器推送 (Server Push)**
- **流优先级 (Stream Priority)**

### HTTP/3 (2022) - RFC 9114
- **基于 QUIC**: UDP 传输
- **改进的拥塞控制**
- **更好的安全性**: TLS 1.3 强制
- **0-RTT 握手**

## HTTP 消息结构

### 请求消息格式

```
Method SP Request-URI SP HTTP-Version CRLF
Header-Field: value CRLF
...
Header-Field: value CRLF
CRLF
[Message Body]
```

**示例：**
```http
POST /api/users HTTP/1.1
Host: api.example.com
Content-Type: application/json
Content-Length: 45
Authorization: Bearer token123

{"name": "John Doe", "email": "john@example.com"}
```

### 响应消息格式

```
HTTP-Version SP Status-Code SP Reason-Phrase CRLF
Header-Field: value CRLF
...
Header-Field: value CRLF
CRLF
[Message Body]
```

**示例：**
```http
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 67
Cache-Control: no-cache

{"id": 123, "name": "John Doe", "email": "john@example.com"}
```

## HTTP 方法

### 安全方法 (Safe Methods)
不会改变服务器状态：

- **GET**: 获取资源
- **HEAD**: 获取响应头（不含主体）
- **OPTIONS**: 获取支持的方法

### 幂等方法 (Idempotent Methods)
多次调用结果相同：

- **PUT**: 更新/创建资源
- **DELETE**: 删除资源
- 所有安全方法

### 主要方法详解

#### GET
```http
GET /users/123 HTTP/1.1
Host: api.example.com
Accept: application/json
```

#### POST
```http
POST /users HTTP/1.1
Host: api.example.com
Content-Type: application/json
Content-Length: 45

{"name": "John", "email": "john@example.com"}
```

#### PUT
```http
PUT /users/123 HTTP/1.1
Host: api.example.com
Content-Type: application/json
Content-Length: 58

{"id": 123, "name": "John Doe", "email": "john@example.com"}
```

#### PATCH
```http
PATCH /users/123 HTTP/1.1
Host: api.example.com
Content-Type: application/json
Content-Length: 25

{"email": "john.doe@example.com"}
```

#### DELETE
```http
DELETE /users/123 HTTP/1.1
Host: api.example.com
```

## 状态码分类

### 1xx: 信息响应
- **100 Continue**: 继续发送请求主体
- **101 Switching Protocols**: 协议切换

### 2xx: 成功响应
- **200 OK**: 请求成功
- **201 Created**: 资源已创建
- **202 Accepted**: 请求已接受，异步处理
- **204 No Content**: 无内容返回

### 3xx: 重定向
- **301 Moved Permanently**: 永久重定向
- **302 Found**: 临时重定向
- **303 See Other**: 查看其他位置
- **304 Not Modified**: 资源未修改（缓存）

### 4xx: 客户端错误
- **400 Bad Request**: 请求语法错误
- **401 Unauthorized**: 未授权
- **403 Forbidden**: 禁止访问
- **404 Not Found**: 资源不存在
- **405 Method Not Allowed**: 方法不允许
- **409 Conflict**: 资源冲突
- **422 Unprocessable Entity**: 语义错误

### 5xx: 服务器错误
- **500 Internal Server Error**: 服务器内部错误
- **501 Not Implemented**: 未实现
- **502 Bad Gateway**: 网关错误
- **503 Service Unavailable**: 服务不可用
- **504 Gateway Timeout**: 网关超时

## HTTP 头信息

### 请求头

#### 通用头
- **Host**: 请求主机 (`api.example.com:8080`)
- **User-Agent**: 客户端信息 (`Mozilla/5.0 ...`)
- **Accept**: 接受的内容类型 (`application/json`)
- **Accept-Language**: 接受的语言 (`zh-CN,zh;q=0.9`)

#### 条件请求
- **If-Modified-Since**: 条件获取 (`Wed, 21 Oct 2015 07:28:00 GMT`)
- **If-None-Match**: ETag 条件 (`"etag-value"`)

#### 认证头
- **Authorization**: 认证信息 (`Bearer token123`)
- **Cookie**: 会话信息 (`session=abc123`)

#### 内容协商
- **Accept**: 媒体类型 (`text/html,application/json`)
- **Accept-Encoding**: 编码 (`gzip, deflate`)
- **Accept-Charset**: 字符集 (`utf-8, iso-8859-1`)

### 响应头

#### 内容信息
- **Content-Type**: 内容类型 (`application/json; charset=utf-8`)
- **Content-Length**: 内容长度 (`1234`)
- **Content-Encoding**: 内容编码 (`gzip`)

#### 缓存控制
- **Cache-Control**: 缓存指令 (`no-cache, max-age=3600`)
- **ETag**: 实体标签 (`"etag-value"`)
- **Last-Modified**: 最后修改时间

#### 重定向
- **Location**: 重定向地址 (`https://example.com/new-url`)

#### 安全
- **Set-Cookie**: 设置 Cookie
- **Strict-Transport-Security**: HSTS
- **X-Frame-Options**: 点击劫持防护

## 内容协商

### 媒体类型 (MIME Types)

**语法：** `type/subtype; parameter=value`

**示例：**
- `text/html`
- `application/json; charset=utf-8`
- `image/png`
- `application/x-git-receive-pack-advertisement`

### 内容编码

- **gzip**: GNU zip 压缩
- **deflate**: zlib 压缩
- **br**: Brotli 压缩

### 字符集

- **UTF-8**: Unicode 编码
- **ISO-8859-1**: Latin-1 编码

## 缓存机制

### 缓存控制头

**请求头：**
```http
Cache-Control: no-cache
If-None-Match: "etag-value"
```

**响应头：**
```http
Cache-Control: max-age=3600, public
ETag: "etag-value"
Last-Modified: Wed, 21 Oct 2015 07:28:00 GMT
```

### 缓存验证

**ETag**: 实体标签，资源版本标识
**Last-Modified**: 最后修改时间

## Cookie 与会话管理

### Cookie 头

**设置 Cookie：**
```http
Set-Cookie: session=abc123; Path=/; HttpOnly; Secure; SameSite=Strict
```

**发送 Cookie：**
```http
Cookie: session=abc123; user_pref=dark_mode
```

### 会话管理

- **基于 Cookie**: 服务器生成 session ID，存储在客户端
- **基于 Token**: JWT 或其他 token 机制
- **无状态**: 每次请求包含完整认证信息

## 连接管理

### HTTP/1.1 持久连接

```http
Connection: keep-alive
```

### HTTP/2 多路复用

- 单个 TCP 连接多个并发流
- 每个流有唯一 ID
- 流优先级和依赖关系

### 连接池

- 复用 TCP 连接
- 减少握手开销
- 自动管理连接生命周期

## 安全考虑

### HTTPS

- **TLS 加密**: 保护传输数据
- **证书验证**: 确保服务器身份
- **HSTS**: 强制 HTTPS

### 常见攻击防护

- **XSS**: `X-XSS-Protection`, `Content-Security-Policy`
- **CSRF**: `SameSite` Cookie 属性
- **点击劫持**: `X-Frame-Options`

## HTTP 在 Git 中的应用

### Git Smart HTTP

**服务发现：**
```http
GET /repo.git/info/refs?service=git-upload-pack HTTP/1.1
Accept: */*
```

**数据传输：**
```http
POST /repo.git/git-upload-pack HTTP/1.1
Content-Type: application/x-git-upload-pack-request
```

### 自定义 MIME 类型

- `application/x-git-receive-pack-advertisement`
- `application/x-git-upload-pack-request`
- `application/x-git-receive-pack-result`

## 调试与开发工具

### 命令行工具

```bash
# curl 基本请求
curl -v https://api.example.com/users

# 带头的请求
curl -H "Authorization: Bearer token" https://api.example.com/users

# POST 请求
curl -X POST -H "Content-Type: application/json" \
  -d '{"name": "John"}' https://api.example.com/users
```

### 浏览器开发者工具

- **Network 面板**: 查看请求响应
- **Headers**: 检查头信息
- **Response**: 查看响应内容

### 代理工具

- **Charles Proxy**
- **Fiddler**
- **mitmproxy**

## 性能优化

### 压缩

```http
Accept-Encoding: gzip, deflate, br
Content-Encoding: gzip
```

### 缓存策略

- **静态资源**: 长缓存时间
- **动态内容**: 条件请求
- **CDN**: 地理分布缓存

### 连接优化

- **HTTP/2**: 多路复用
- **连接池**: 复用连接
- **域名分片**: 并行下载

## 参考标准

- **RFC 7230**: HTTP/1.1 消息语法和路由
- **RFC 7231**: HTTP/1.1 语义和内容
- **RFC 7232**: HTTP 条件请求
- **RFC 7233**: HTTP 范围请求
- **RFC 7234**: HTTP 缓存
- **RFC 7235**: HTTP 认证
- **RFC 7540**: HTTP/2
- **RFC 9114**: HTTP/3

---

本文档系统介绍了 HTTP 协议的核心概念、详细规范和实际应用，是构建 Web 服务和理解现代互联网架构的基础知识。
