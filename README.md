# aq3cms

<div align="center">

# aq3cms

**基于 Go 语言的现代化内容管理系统**

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://docker.com)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen)](https://github.com/your-repo/aq3cms)

</div>

## 📖 项目介绍

aq3cms 是一个基于 Go 语言开发的现代化内容管理系统，完全兼容 aq3CMS 数据库结构和模板语法。它结合了 Go 语言的高性能特性和传统 CMS 的易用性，为用户提供了一个安全、高效、易扩展的内容管理解决方案。

### 🎯 设计目标

- **性能优先**：充分利用 Go 语言的并发特性，提供卓越的性能表现
- **完全兼容**：100% 兼容 aq3CMS 数据库结构和模板语法，无缝迁移
- **现代化架构**：采用微服务架构设计，支持容器化部署
- **安全可靠**：内置多层安全防护，保障系统和数据安全
- **易于扩展**：插件化架构，支持功能扩展和定制开发

## ✨ 核心特性

### 🚀 高性能架构
- **并发处理**：基于 Go 语言的 goroutine，支持高并发访问
- **内存优化**：高效的内存管理，降低资源消耗
- **缓存系统**：多级缓存策略，显著提升响应速度
- **静态化**：支持页面静态化，减轻服务器压力

### 🔄 完美兼容
- **数据库兼容**：完全兼容 aq3CMS 数据库结构
- **模板兼容**：支持 aq3CMS 模板标签语法
- **URL兼容**：保持与 aq3CMS 相同的 URL 结构
- **功能兼容**：核心功能与 aq3CMS 保持一致

### 🛡️ 安全防护
- **密码安全**：使用 Argon2id 加密算法
- **XSS 防护**：输入过滤和输出编码
- **CSRF 防护**：令牌验证机制
- **SQL 注入防护**：参数化查询
- **文件上传安全**：严格的文件类型检查

### 📱 现代化体验
- **响应式设计**：完美支持 PC 和移动端
- **RESTful API**：标准的 API 接口设计
- **前后端分离**：支持前后端分离架构
- **实时更新**：WebSocket 支持实时功能

### 🌐 国际化支持
- **多语言**：支持多语言界面
- **本地化**：支持不同地区的本地化设置
- **时区支持**：自动处理时区转换
- **字符编码**：完整的 UTF-8 支持

### 🔌 插件生态
- **插件系统**：灵活的插件架构
- **钩子机制**：丰富的扩展点
- **热插拔**：支持插件的动态加载和卸载
- **API 扩展**：插件可扩展 API 接口

## 🏗️ 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Nginx/CDN     │    │   Load Balancer │    │   Monitoring    │
│   (反向代理)     │    │   (负载均衡)     │    │   (监控告警)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   aq3cms App 1   │    │   aq3cms App 2   │    │   aq3cms App N   │
│   (应用实例)     │    │   (应用实例)     │    │   (应用实例)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   MySQL Master │    │   Redis Cluster │    │   File Storage  │
│   (主数据库)     │    │   (缓存集群)     │    │   (文件存储)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │
┌─────────────────┐
│   MySQL Slave   │
│   (从数据库)     │
└─────────────────┘
```

### 目录结构

```
aq3cms/
├── cmd/                  # 命令行工具
│   └── server/           # 服务器启动程序
├── config/               # 配置文件
├── internal/             # 内部包
│   ├── controller/       # 控制器
│   │   ├── admin/        # 后台控制器
│   │   ├── api/          # API控制器
│   │   └── frontend/     # 前台控制器
│   ├── middleware/       # 中间件
│   ├── model/            # 数据模型
│   ├── service/          # 业务逻辑
│   └── template/         # 模板引擎
│       └── tags/         # 模板标签
├── pkg/                  # 可重用包
│   ├── cache/            # 缓存系统
│   ├── database/         # 数据库操作
│   ├── logger/           # 日志系统
│   ├── plugin/           # 插件系统
│   ├── security/         # 安全相关
│   └── i18n/             # 国际化
├── docker/               # Docker 配置
│   ├── nginx/            # Nginx 配置
│   ├── mysql/            # MySQL 配置
│   └── redis/            # Redis 配置
├── static/               # 静态资源
├── templets/             # 模板文件
├── uploads/              # 上传文件
├── logs/                 # 日志文件
└── data/                 # 数据文件
```

## 🚀 快速开始

### 方式一：Docker Compose（推荐）

```bash
# 1. 克隆项目
git clone https://github.com/your-repo/aq3cms.git
cd aq3cms

# 2. 启动服务
docker-compose up -d

# 3. 访问系统
# 前台：http://localhost
# 后台：http://localhost/admin
# 默认账号：admin / admin123
```

### 方式二：本地构建

```bash
# 1. 安装 Go 1.23+
# 2. 安装 MySQL 8.0+
# 3. 安装 Redis 7.0+（可选）

# 4. 克隆项目
git clone https://github.com/your-repo/aq3cms.git
cd aq3cms

# 5. 配置数据库
# 编辑 config.yaml 文件

# 6. 构建应用
make build

# 7. 运行应用
./bin/aq3cms
```

## 📋 系统要求

### 最低要求
- **CPU**: 1 核心
- **内存**: 512MB
- **存储**: 1GB
- **Go**: 1.21+
- **MySQL**: 5.7+
- **Redis**: 6.0+（可选）

### 推荐配置
- **CPU**: 2 核心以上
- **内存**: 2GB 以上
- **存储**: 10GB 以上
- **Go**: 1.23+
- **MySQL**: 8.0+
- **Redis**: 7.0+

### 支持的操作系统
- Linux (Ubuntu 20.04+, CentOS 8+, Debian 11+)
- macOS (10.15+)
- Windows (10+)

## 📚 文档导航

### 用户文档
- [安装指南](docs/installation.md)
- [用户手册](docs/user-guide.md)
- [模板开发](docs/template.md)
- [插件开发](docs/plugin.md)

### 开发文档
- [开发指南](docs/development.md)
- [API 文档](docs/api.md)
- [数据库设计](docs/database.md)
- [架构设计](docs/architecture.md)

### 运维文档
- [部署文档](DEPLOYMENT.md)
- [监控告警](docs/monitoring.md)
- [性能调优](docs/performance.md)
- [故障排除](docs/troubleshooting.md)

## 🛠️ 开发工具

### 必需工具
- [Go](https://golang.org/) - 编程语言
- [Docker](https://docker.com/) - 容器化
- [Make](https://www.gnu.org/software/make/) - 构建工具

### 推荐工具
- [Air](https://github.com/cosmtrek/air) - 热重载
- [Delve](https://github.com/go-delve/delve) - 调试工具
- [golangci-lint](https://golangci-lint.run/) - 代码检查

### IDE 推荐
- [VS Code](https://code.visualstudio.com/) + Go 扩展
- [GoLand](https://www.jetbrains.com/go/)
- [Vim](https://www.vim.org/) + vim-go

## 🤝 贡献指南

我们欢迎所有形式的贡献！请阅读 [贡献指南](CONTRIBUTING.md) 了解详情。

### 贡献方式
- 🐛 报告 Bug
- 💡 提出新功能建议
- 📝 改进文档
- 🔧 提交代码修复
- 🌟 推广项目

### 开发流程
1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 📄 许可证

本项目采用 [Apache-2.0 license 许可证](LICENSE)。

## 🙏 致谢

感谢以下开源项目的支持：

- [Go](https://golang.org/) - 优秀的编程语言
- [Gorilla](https://www.gorillatoolkit.org/) - Web 工具包
- [GORM](https://gorm.io/) - ORM 库
- [Redis](https://redis.io/) - 内存数据库
- [MySQL](https://mysql.com/) - 关系型数据库

## 📞 联系我们

- **官网**: https://aq3.cn
- **文档**: https://docs.aq3.cn
- **社区**: https://community.aq3.cn
- **邮箱**: support@aq3.cn

---

<div align="center">

**如果这个项目对您有帮助，请给我们一个 ⭐️**

Made with ❤️ by aq3cms Team

</div>
