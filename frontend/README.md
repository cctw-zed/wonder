# Wonder Frontend

现代化的前端Web应用程序，使用Next.js、TypeScript和Tailwind CSS构建。

## 🚀 技术栈

- **框架**: Next.js 15.5 (React 19)
- **语言**: TypeScript
- **样式**: Tailwind CSS v4
- **UI组件**: shadcn/ui
- **状态管理**: Zustand
- **表单处理**: React Hook Form + Zod
- **HTTP客户端**: Axios
- **构建工具**: Next.js (内置Webpack)

## 📋 功能特性

### ✅ 已实现功能
- 🎨 **现代化UI设计**: 使用shadcn/ui组件和Tailwind CSS
- 🔐 **用户认证系统**: 完整的登录/注册流程
- 📱 **响应式设计**: 支持移动端和桌面端
- 🔄 **状态管理**: 基于Zustand的认证状态管理
- ✅ **表单验证**: 使用Zod进行强类型验证
- 🔒 **TypeScript支持**: 完整的类型安全
- 🐳 **Docker支持**: 容器化部署配置

### 📄 页面结构
- **首页** (`/`) - 产品介绍和导航
- **登录页面** (`/login`) - 用户登录
- **注册页面** (`/register`) - 用户注册
- **仪表板** (`/dashboard`) - 用户个人中心

## 🛠️ 开发指南

### 环境要求
- Node.js 18+
- npm 或 yarn

### 快速开始

```bash
# 安装依赖
npm install

# 复制环境变量文件
cp .env.example .env.local

# 启动开发服务器
npm run dev
```

应用将在 `http://localhost:3001` 启动。

### 可用命令

```bash
# 开发服务器
npm run dev

# 生产构建
npm run build

# 启动生产服务器
npm start

# 代码检查
npm run lint
```

## 🔧 配置

### 环境变量

在 `.env.local` 文件中配置以下变量：

```bash
# 后端API地址
NEXT_PUBLIC_API_URL=http://localhost:8080

# 应用配置
NEXT_PUBLIC_APP_NAME=Wonder
NEXT_PUBLIC_APP_VERSION=1.0.0
```

### API集成

前端通过以下API端点与后端通信：

### 认证相关
- `POST /api/v1/users/register` - 用户注册 (公开)
- `POST /api/v1/auth/login` - 用户登录 (公开)
- `POST /api/v1/auth/logout` - 用户登出 (需认证)
- `GET /api/v1/auth/me` - 获取当前用户信息 (需认证)

### 用户管理
- `GET /api/v1/users` - 获取用户列表 (可选认证)
- `GET /api/v1/users/:id` - 通过ID获取用户信息 (需认证)
- `PUT /api/v1/users/:id` - 更新用户信息 (需认证)
- `DELETE /api/v1/users/:id` - 删除用户 (需认证)

## 📁 项目结构

```
frontend/
├── src/
│   ├── app/                    # Next.js App Router页面
│   │   ├── page.tsx           # 首页
│   │   ├── login/             # 登录页面
│   │   ├── register/          # 注册页面
│   │   └── dashboard/         # 仪表板
│   ├── components/
│   │   └── ui/                # shadcn/ui组件
│   ├── lib/
│   │   ├── api.ts            # API客户端配置
│   │   └── utils.ts          # 工具函数
│   └── store/
│       └── auth.ts           # 认证状态管理
├── public/                    # 静态资源
├── Dockerfile                # Docker构建配置
├── next.config.ts            # Next.js配置
├── tailwind.config.ts        # Tailwind CSS配置
└── components.json           # shadcn/ui配置
```

## 🎨 UI/UX设计

### 设计系统
- **主色调**: 蓝色系 (#3B82F6)
- **字体**: 系统默认字体栈
- **响应式断点**: 遵循Tailwind CSS标准
- **组件**: 基于shadcn/ui的现代化组件

### 用户体验
- 清晰的导航结构
- 直观的表单验证
- 流畅的页面转换
- 友好的错误处理

## 🔐 认证流程

### 登录流程
1. 用户输入邮箱和密码
2. 前端验证表单数据
3. 发送API请求到后端
4. 后端返回JWT Token
5. 前端保存Token并跳转到仪表板

### 注册流程
1. 用户填写注册信息
2. 前端验证表单数据（包括密码确认）
3. 发送API请求到后端
4. 注册成功后自动登录
5. 跳转到仪表板

### 状态管理
使用Zustand管理认证状态：
- 用户信息存储
- Token管理
- 自动登录/登出
- 错误状态处理

## 🐳 Docker部署

### 构建镜像

```bash
# 构建Docker镜像
docker build -t wonder-frontend .

# 运行容器
docker run -p 3001:3001 wonder-frontend
```

### 多阶段构建
Dockerfile使用多阶段构建：
1. **构建阶段**: 安装依赖并构建应用
2. **运行阶段**: 创建轻量级生产镜像

### 环境变量
Docker部署时可通过环境变量配置：

```bash
docker run -p 3001:3001 \
  -e NEXT_PUBLIC_API_URL=http://your-backend-url \
  wonder-frontend
```

## 🔄 与后端集成

### API通信
- 使用Axios作为HTTP客户端
- 自动添加认证Token到请求头
- 统一的错误处理
- 请求/响应拦截器

### 数据流
```
前端组件 → Zustand Store → API Service → 后端
         ←                ←            ←
```

## 🧪 测试（规划中）

### 测试策略
- **单元测试**: 组件逻辑测试
- **集成测试**: API集成测试
- **E2E测试**: 完整用户流程测试

### 工具选择
- Jest + React Testing Library
- Cypress 或 Playwright (E2E)

## 📈 性能优化

### 已实现优化
- **代码分割**: Next.js自动代码分割
- **图片优化**: Next.js Image组件
- **缓存策略**: 适当的缓存配置
- **Tree Shaking**: 自动删除未使用代码

### 构建优化
- Standalone输出模式（适合Docker）
- 生产环境优化配置
- 压缩和Minification

## 🚀 部署

### 开发环境
```bash
npm run dev
```

### 生产环境
```bash
npm run build
npm start
```

### Docker容器
```bash
docker-compose up -d
```

## 📚 学习资源

对于前端新手，推荐学习资源：

### 基础知识
- [React官方文档](https://react.dev/)
- [Next.js文档](https://nextjs.org/docs)
- [TypeScript手册](https://www.typescriptlang.org/docs/)

### UI/样式
- [Tailwind CSS文档](https://tailwindcss.com/docs)
- [shadcn/ui组件库](https://ui.shadcn.com/)

### 状态管理
- [Zustand文档](https://docs.pmnd.rs/zustand/getting-started/introduction)

## 🤝 贡献指南

1. Fork项目
2. 创建功能分支
3. 编写代码和测试
4. 提交Pull Request

## 📄 许可证

MIT License - 详见 LICENSE 文件

---

**Wonder Frontend** - 现代化、类型安全的React应用程序
