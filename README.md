# Hyper Pen - 在线笔记应用

Hyper Pen 是一个基于Vue3和Go的在线笔记应用，支持Markdown格式的笔记编辑和实时预览。

## 功能特点

- Markdown格式的笔记编辑和实时预览
- 用户注册和登录
- GitHub和微信第三方登录
- 笔记的创建、编辑、删除和查看

## 技术栈

### 前端
- Vue 3
- Element Plus
- Vite
- Marked (Markdown解析)
- Highlight.js (代码高亮)

### 后端
- Go
- Iris Web框架
- GORM
- SQLite

## 项目结构

```
.
├── hyper-pen-ui/          # 前端项目
│   ├── src/               # 源代码
│   │   ├── views/         # 页面组件
│   │   ├── router/        # 路由配置
│   │   └── main.js        # 入口文件
│   └── package.json       # 依赖配置
│
└── hyper-pen-service/     # 后端项目
    ├── handlers/          # 请求处理
    ├── models/            # 数据模型
    ├── db/                # 数据库相关
    └── main.go            # 入口文件
```

## 开发环境设置

### 前端开发

1. 进入前端项目目录：
```bash
cd hyper-pen-ui
```

2. 安装依赖：
```bash
npm install
```

3. 启动开发服务器：
```bash
npm run dev
```

### 后端开发

1. 进入后端项目目录：
```bash
cd hyper-pen-service
```

2. 安装依赖：
```bash
go mod tidy
```

3. 启动服务器：
```bash
go run main.go
```

## API文档

### 认证相关

- POST /api/auth/login - 用户登录
- POST /api/auth/register - 用户注册

### 笔记相关

- GET /api/notes - 获取笔记列表
- POST /api/notes - 创建新笔记
- PUT /api/notes/:id - 更新笔记
- DELETE /api/notes/:id - 删除笔记

## 待实现功能

- [ ] JWT认证
- [ ] GitHub OAuth登录
- [ ] 微信扫码登录
- [ ] 笔记分类和标签
- [ ] 笔记搜索功能
- [ ] 笔记分享功能 