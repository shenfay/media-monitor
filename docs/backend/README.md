# 项目文档导航

## 📚 文档分类

### 🚀 快速开始
- [快速开始指南](development/GETTING_STARTED.md) - 5 分钟运行项目
- [开发指南](development/DEVELOPMENT_GUIDE.md) - 开发规范与流程
- [代码注释规范](development/CODE_COMMENT_GUIDELINES.md) - 注释规范与示例

### 🏗️ 架构设计
- [DDD 架构设计](architecture/DDD_ARCHITECTURE.md) - 架构理念与分层设计
- [领域模型](architecture/DOMAIN_MODEL.md) - 聚合根、实体、值对象、RBAC 权限领域、操作日志聚合

### 🗄️ 数据库
- [数据库设计](database/SCHEMA_DESIGN.md) - 表结构、ER 图、迁移管理

### 📡 API 文档
- 启动服务后访问：http://localhost:8080/swagger/index.html

### 🐳 部署运维
- [Docker 部署](deployment/DOCKER_DEPLOYMENT.md) - 容器化部署

### 📊 监控运维
- [监控配置指南](operations/MONITORING_SETUP.md) - Prometheus + Grafana 配置
- [故障排查指南](operations/TROUBLESHOOTING.md) - 常见问题诊断与解决

### 🖥️ 管理后台前端
- [前端架构文档](../admin/ARCHITECTURE.md) - 路由、状态管理、i18n、组件体系

### 📝 产品文档
- [产品需求文档](../product/PRD.md)
- [MVP 功能清单与验收标准](../product/MVP功能清单与验收标准.md)
- [用户故事文档](../product/用户故事文档.md)
- [文案内容体系](../product/文案内容体系.md)
- [核心交互流程图](../product/核心交互流程图.md)

### 🎨 UI 设计
- [UI 设计规范](../UI_DESIGN_SPEC.md) - 布局、色彩、字体、按钮、表格等设计规范

## 📖 文档规范

### 格式标准
- ✅ 使用 Markdown 格式
- ✅ 架构图/流程图使用 Mermaid 语法
- ✅ 代码块标注语言类型
- ✅ 表格对齐整齐

### 命名规范
- 文件名：大写字母 + 下划线（如 `GETTING_STARTED.md`）
- 目录名：小写字母（如 `architecture`、`development`）

### 内容结构
每个文档应包含：
1. **标题** - 清晰的文档主题
2. **概述** - 文档目的和适用范围
3. **正文** - 结构化内容（使用标题分级）
4. **示例** - 代码示例/图表/配置
5. **参考** - 相关链接和延伸阅读

## 🎯 文档优先级

### P0 - 核心文档（必须）
- [x] README.md（项目根目录）
- [x] 快速开始指南
- [x] DDD 架构设计
- [x] 代码注释规范

### P1 - 重要文档（已完成）
- [x] 领域模型（含 RBAC 权限领域、操作日志聚合）
- [x] 数据库设计（含 RBAC 表、菜单表、统一操作日志表）
- [x] 开发指南
- [x] Docker 部署
- [x] 监控配置
- [x] 故障排查
- [x] 管理后台前端架构文档

### P2 - 辅助文档（按需补充）
- [ ] 事件风暴
- [ ] 迁移指南
- [ ] 生产检查清单
- [ ] 测试指南
- [ ] 性能调优

## 🤝 贡献指南

### 添加新文档
1. 确定文档分类（architecture/development/database/deployment/operations）
2. 使用大写文件名（如 `NEW_FEATURE.md`）
3. 在本文档中添加链接
4. 提交 Pull Request

### 更新文档
1. 保持内容准确、最新
2. 更新文档修改日期
3. 如有重大变更，通知团队成员

## 📅 文档维护

- **代码变更时**：同步更新相关文档
- **每月审查**：检查文档准确性
- **版本发布时**：更新版本号和变更记录

---

**最后更新**：2026-07-09  
**文档状态**：✅ 核心文档已完成，持续优化中
