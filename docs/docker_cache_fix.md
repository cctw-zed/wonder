🔧 问题原因

错误: failed to compute cache key: short read: expected 57323671 bytes but got 43388928: unexpected EOF

这是典型的 Docker 构建缓存损坏问题，通常由以下原因引起：
- Docker 构建过程中被中断
- 网络问题导致镜像下载不完整
- 磁盘空间不足导致缓存文件损坏
- Docker 版本升级后缓存不兼容

🛠️ 解决方案

1. 创建了专用修复脚本

文件: scripts/fix-docker-cache.sh

脚本自动执行以下操作：
- 🛑 停止所有 Wonder 服务
- 🧹 清理 Docker 系统缓存
- 🗑️ 移除损坏的镜像
- 🧽 清理构建缓存

2. 集成到 Makefile

添加了新命令：
make fix-docker-cache

3. 使用方法

当遇到类似错误时：
# 快速修复
make fix-docker-cache

# 然后重新构建
make setup

或者手动步骤：
# 1. 清理缓存
docker system prune -a -f
docker builder prune -a -f

# 2. 停止服务
make stop-all

# 3. 无缓存重建
make rebuild-no-cache

📋 更新的工作流命令

现在你有了完整的错误处理工具集：

# 🔧 错误修复
make fix-docker-cache     # 修复 Docker 缓存问题

# 🚀 环境管理
make setup                # 完整环境设置
make rebuild-no-cache     # 无缓存重建
make stop-all             # 停止所有服务
make clean-all           # 完全重置（删除数据）

# 📊 状态检查
make status              # 查看服务状态
make urls               # 查看服务地址

🎯 预防措施

为避免将来再次出现缓存问题：

1. 定期清理: 每周运行 docker system prune -f
2. 监控磁盘空间: 确保 Docker 目录有足够空间
3. 使用无缓存构建: 重要更改时使用 make rebuild-no-cache
4. 避免中断构建: 让构建过程完整完成

💡 故障排除指南

如果再次遇到构建问题：

1. 首先尝试: make fix-docker-cache
2. 检查磁盘空间: df -h
3. 检查 Docker 状态: docker system df
4. 如果问题持续: 重启 Docker 服务

⏺ 现在你有了一个强大的工具来处理 Docker 缓存问题！当你准备好时，可以运行：

make fix-docker-cache
make setup

这将给你一个干净的环境开始。第一次构建可能需要一些时间来下载镜像，但之后的构建会很快。🚀