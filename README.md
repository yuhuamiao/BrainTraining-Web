# 脑力训练台

一个可直接运行和部署的全栈大脑训练应用。前端使用 Vue 3，后端使用 Go + Gin，默认由 SQLite 持久化用户与成绩；也可通过环境变量切换到 MySQL。

## 已实现功能

- 用户注册、登录、JWT 鉴权和用户名修改
- 舒尔特方格：30 秒顺序视觉搜索
- 多色文字：WebSocket 连续颜色刺激，断线自动降级
- 瞬时记忆：三档难度、10 轮空间位置记忆
- 公交人数：三档速度和轮次的工作记忆心算
- 数独挑战：三档可解题目、错误与用时记录
- 个人成绩中心：按训练筛选记录、总次数和平均准确率
- SQLite 开箱即用、可选 MySQL、单容器 Docker 部署
- 后端 API 回归测试与前端核心游戏逻辑测试

## 技术栈

| 层 | 技术 |
| --- | --- |
| 前端 | Vue 3、Vue Router、Pinia、Axios、Vite |
| 游戏 | sudoku-gen、WebSocket |
| 后端 | Go 1.24、Gin、GORM、JWT |
| 数据库 | SQLite（默认）或 MySQL |
| 部署 | Docker、Docker Compose |

## 本地运行

环境要求：Node.js 22、Go 1.24。无需提前安装数据库。

1. 启动后端：

   ```powershell
   cd backend
   Copy-Item .env.example .env
   go run .
   ```

2. 另开终端启动前端：

   ```powershell
   cd frontend
   npm install
   npm run dev
   ```

3. 打开 `http://localhost:5173`，注册账号后即可训练。Vite 会代理 `/api` 和 WebSocket 到 `http://127.0.0.1:8000`。

后端会自动创建 `backend/data/brain-training.db` 并迁移数据表。API 健康检查位于 `http://localhost:8000/health`，Swagger 页面位于 `http://localhost:8000/swagger/index.html`。

## Docker 部署

生产部署必须先设置至少 32 个字符的随机 JWT 密钥，未设置时 Compose 会直接拒绝启动：

```bash
export JWT_SECRET="$(openssl rand -hex 32)"
docker compose up -d --build
```

打开 `http://localhost:8000`。前后端打包在同一容器中，数据库和头像分别保存在 Docker volume，重建容器不会丢失。

公网部署时应在容器前使用 Nginx、Caddy 等反向代理启用 HTTPS，并对登录、注册和头像上传接口配置请求频率限制；不要直接把开发模式服务暴露到公网。

Docker 构建默认通过 `https://goproxy.cn,direct` 下载 Go 模块；需要使用其他代理时，可在构建前设置 `GOPROXY` 环境变量覆盖。

## 配置

后端会读取 `backend/.env` 或进程环境变量。

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `PORT` | `8000` | HTTP 服务端口 |
| `JWT_SECRET` | 仅 debug 模式有本地开发值 | release/Docker 模式必须设置至少 32 个字符的随机密钥 |
| `DB_DRIVER` | `sqlite` | `sqlite` 或 `mysql` |
| `DB_PATH` | `./data/brain-training.db` | SQLite 文件路径；默认启用 WAL、外键和 5 秒锁等待 |
| `DB_HOST` / `DB_PORT` | `127.0.0.1` / `3306` | MySQL 地址 |
| `DB_USER` / `DB_PASSWORD` / `DB_NAME` | - | MySQL 凭据与库名 |
| `CORS_ALLOWED_ORIGIN` | 空 | 跨域开发时允许的精确来源 |
| `FRONTEND_DIR` | `./static` | Gin 托管的前端构建目录 |

## API

公开接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `POST` | `/api/v1/register` | 注册 |
| `POST` | `/api/v1/login` | 登录并获取 JWT |
| `GET` | `/api/v1/schulte/matrix` | 舒尔特矩阵 |
| `GET` | `/api/v1/color_words/matrix` | 多色文字矩阵 |
| `GET (WS)` | `/api/v1/color_words/color` | 颜色流 |
| `GET` | `/api/v1/memory/matrix` | 记忆方格 |

以下接口需要 `Authorization: Bearer <token>`：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `POST` | `/api/v1/{schulte,color_words,memory,bus,sudoku}/scores` | 保存训练成绩 |
| `GET` | `/api/v1/user/scores` | 当前用户全部成绩 |
| `GET` / `PATCH` | `/api/v1/user/me` | 获取或修改当前用户 |
| `POST` | `/api/v1/user/photo` | 上传头像（最大 2 MB） |

受保护接口的用户身份只从 JWT 读取，请求体或查询参数中的 `userId` 不会改变数据归属。

## 测试与构建

```powershell
cd backend
go test ./...

cd ../frontend
npm test
npm run build
```

后端测试覆盖注册登录、鉴权、成绩用户归属和训练数据校验；前端测试覆盖公交人数生成不变量和数独生成器集成。

## 目录结构

```text
BrainTraining-Web/
├── backend/           # Gin API、GORM 模型、SQLite/MySQL 配置与测试
├── frontend/          # Vue 训练页面、成绩中心和前端逻辑测试
├── Dockerfile         # 前后端多阶段构建
├── docker-compose.yml # 单服务部署与持久化卷
└── README.md
```
