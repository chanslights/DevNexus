# Git Smart HTTP Protocol 详解

## 概述

Git Smart HTTP 协议是 Git 支持的一种传输协议，允许通过 HTTP/HTTPS 进行高效的版本控制操作。它是 GitHub、GitLab 等平台的核心基础设施，支持 `git clone`、`git fetch`、`git pull` 和 `git push` 操作。

## 协议历史

- **Dumb HTTP Protocol**: Git 早期版本使用，支持只读操作，通过多次 HTTP GET 请求获取松散对象
- **Smart HTTP Protocol**: Git 1.6.6 引入，支持双向操作，使用高效的 packfile 传输

## 协议架构

### 两阶段设计

Git Smart HTTP 协议采用两阶段设计：

1. **服务发现阶段 (Discovery Phase)**: 客户端发现可用服务和引用
2. **数据传输阶段 (Data Transfer Phase)**: 实际的数据传输

### URL 结构

```
http[s]://host[:port]/path/to/repo.git/action/parameters
```

**示例：**
- `https://github.com/user/repo.git/info/refs`
- `https://github.com/user/repo.git/git-upload-pack`

## 详细协议流程

### 阶段1：服务发现 (Discovery)

#### 1.1 客户端请求

**Push 操作：**
```http
GET /repo.git/info/refs?service=git-receive-pack HTTP/1.1
Host: example.com
User-Agent: git/2.34.1
Accept: */*
```

**Pull 操作：**
```http
GET /repo.git/info/refs?service=git-upload-pack HTTP/1.1
Host: example.com
User-Agent: git/2.34.1
Accept: */*
```

#### 1.2 服务端响应

**响应头：**
```http
HTTP/1.1 200 OK
Content-Type: application/x-git-receive-pack-advertisement
Cache-Control: no-cache
```

**响应体格式：**
```
001e# service=git-receive-pack\n
0000
005b0000000000000000000000000000000000000000 refs/heads/main\n
003f0000000000000000000000000000000000000000 refs/heads/feature\n
0000
```

**PKT-LINE 格式：**
- 前4个字符：十六进制长度（包括自身4字节）
- 内容：实际数据
- `0000`：结束标记

### 阶段2：数据传输 (RPC)

#### 2.1 Push 操作

**客户端请求：**
```http
POST /repo.git/git-receive-pack HTTP/1.1
Host: example.com
User-Agent: git/2.34.1
Content-Type: application/x-git-receive-pack-request
Accept: application/x-git-receive-pack-result
Content-Length: 1234

[packfile data]
```

**服务端响应：**
```http
HTTP/1.1 200 OK
Content-Type: application/x-git-receive-pack-result

[unpack status]
[ref updates]
```

#### 2.2 Pull 操作

**客户端请求：**
```http
POST /repo.git/git-upload-pack HTTP/1.1
Host: example.com
Content-Type: application/x-git-upload-pack-request
Accept: application/x-git-upload-pack-result

[want/have lines]
0000
[packfile request]
```

## MIME 类型规范

| 操作 | 请求 Content-Type | 响应 Content-Type |
|------|------------------|------------------|
| Push Discovery | - | `application/x-git-receive-pack-advertisement` |
| Push Data | `application/x-git-receive-pack-request` | `application/x-git-receive-pack-result` |
| Pull Discovery | - | `application/x-git-upload-pack-advertisement` |
| Pull Data | `application/x-git-upload-pack-request` | `application/x-git-upload-pack-result` |

## Packfile 格式

### Packfile 结构

```
[signature]    "PACK"
[version]      4字节大端序版本号
[object count] 4字节大端序对象数量
[objects]      压缩的Git对象
[checksum]     20字节SHA-1校验和
```

### Delta 压缩

- **REF_DELTA**: 基于其他对象的差异
- **OFS_DELTA**: 基于同一packfile中其他对象的差异

## 认证与授权

### HTTP 基本认证

```http
Authorization: Basic <base64-encoded-credentials>
```

### Bearer Token

```http
Authorization: Bearer <token>
```

## 错误处理

### 常见错误响应

**400 Bad Request:**
```http
HTTP/1.1 400 Bad Request
Content-Type: text/plain

Invalid request format
```

**401 Unauthorized:**
```http
HTTP/1.1 401 Unauthorized
WWW-Authenticate: Basic realm="Git Repository"
```

**403 Forbidden:**
```http
HTTP/1.1 403 Forbidden
Content-Type: text/plain

Access denied
```

## 实现注意事项

### 1. 无状态设计

- 每个请求都是独立的
- 使用 `--stateless-rpc` Git 选项
- 不依赖服务器端会话状态

### 2. 并发处理

- 支持多个客户端同时操作
- 需要处理锁竞争（reference updates）
- 原子性保证

### 3. 性能优化

- **Packfile 重用**: 避免重复传输
- **浅克隆 (Shallow Clone)**: `--depth` 参数
- **增量获取**: 使用 `have` 行指定已有对象

### 4. 安全考虑

- **路径遍历防护**: 验证仓库路径
- **资源限制**: 限制 packfile 大小
- **访问控制**: 仓库级别的权限检查

## 实际应用示例

### 使用 curl 测试

```bash
# 发现服务
curl -v "http://localhost:8080/repo.git/info/refs?service=git-upload-pack"

# 模拟push
echo -e "0032want 0000000000000000000000000000000000000000\n0000" | \
curl -v -X POST \
  -H "Content-Type: application/x-git-upload-pack-request" \
  --data-binary @- \
  http://localhost:8080/repo.git/git-upload-pack
```

## 协议扩展

### Git LFS (Large File Storage)

- 将大文件存储在外部服务
- 使用指针文件替代实际内容
- 扩展了 Smart HTTP 协议

### Git Wire Protocol v2

- 更高效的批量操作
- 改进的错误报告
- 更好的可扩展性

## 调试技巧

### 启用 Git 调试

```bash
export GIT_CURL_VERBOSE=1
export GIT_TRACE=1
git clone http://example.com/repo.git
```

### Wireshark 抓包

使用 Wireshark 过滤 `http` 流量分析协议交互。

## 参考资料

- [Git Documentation - HTTP Protocol](https://git-scm.com/docs/http-protocol)
- [Git Packfile Format](https://git-scm.com/docs/pack-format)
- [RFC 7230 - HTTP/1.1 Message Syntax and Routing](https://tools.ietf.org/html/rfc7230)

---

本文档涵盖了 Git Smart HTTP 协议的核心概念、详细流程和实现细节，可作为技术参考和知识库资料使用。
