# course-bot

面向京都大学学生的选课推荐 chatbot。

## 目录结构
- backend/   后端 API（FastAPI）
- rag/  数据抓取与向量化
- frontend/  前端
- docs/      技术说明

## 本地启动
前提：已安装并启动 Docker Desktop。

1. 复制环境变量文件并按需修改：`cp .env.example .env`
2. 启动：`docker compose up --build`
3. 验证：浏览器访问 http://localhost:8000/health ，应返回 {"status":"ok"}
4. 接口文档：http://localhost:8000/docs
5. 停止：Ctrl+C，或另开终端执行 `docker compose down`

## 协作约定
- 不直接推送到 main，从 main 切分支开发，通过 PR 合并
- 分支命名：feature/xxx、fix/xxx、infra/xxx
- PR 保持小而聚焦，至少一位同学 review 后合并
- 合并方式：Squash merge
- 密钥不进仓库，见 .env.example

## 文档
Confluence：https://kyotou.atlassian.net/wiki/spaces/~712020da7a3e93863142ee86d45aecdacdc9e3/pages/edit-v2/360449?draftShareId=ee4217a6-65c4-4752-ac7c-f3b96d878e8a
