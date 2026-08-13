# 后台 UI 设计规范

## 1. 布局规范 (Layout)

### 页面整体结构

```
+-------------------------------------------------------+
| TopBar (50px)                                         |
+--------+----------------------------------------------+
|        |  Page Header (标题 + 操作按钮)                |
| Sidebar|  20px padding top, 28px padding left          |
| (220px |  间距 16px                                    |
| 收起时 |  + Filter Bar (筛选栏)                        |
|  56px) |    0 28px padding, 间距 16px                  |
|        |  + Content (表格 / 卡片)                      |
|        |    0 28px 20px padding                        |
|        |    - 内容区有边框: 1px solid #efeae2, 圆角12px |
+--------+----------------------------------------------+
```

### 尺寸

| 区域 | 属性 | 值 |
|------|------|-----|
| TopBar | 高度 | 50px |
| TopBar | 水平内边距 | 16px 28px |
| Sidebar | 展开宽度 | 224px (220px content + 4px border) |
| Sidebar | 收起宽度 | 56px |
| Page Header | 上内边距 | 20px |
| Page Header | 水平内边距 | 28px |
| Page Header 与下方间距 | height | 16px |
| Filter Bar | 水平内边距 | 28px |
| Filter Bar 与表格间距 | height | 16px |
| Content Area | 内边距 | 0 28px 20px |
| 侧边栏背景色 | --sidebar-bg | #F5F3EF（暖米色） |
| 主内容区背景色 | --main-bg | #FFFFFF |

---

## 2. 字体与排版 (Typography)

### 字号层级

| 用途 | 字号 | 字重 | 颜色 |
|------|------|------|------|
| 页面主标题 (h2) | 20px | 600 (SemiBold) | #2b2b2b |
| 卡片标题 | 14px | 600 | inherit |
| 表格表头 | 12px | 600 | #8a8276, uppercase |
| 表格正文 | 13px | 400 | #2b2b2b |
| 表单标签 | 13px | 500 | inherit |
| 面包屑 | 13px | 400/500 | 上层 #b0a89a, 当前 #6b6258 |
| 按钮文字 | 13px | 500 | 视类型而定 |
| 标签文字 | 12px | 500 | 视状态而定 |
| 辅助/占位文字 | 13px | 400 | #b0a89a / #c4bdb0 |

### 字体栈

```
--font-family: -apple-system, "PingFang SC", "Helvetica Neue", "Microsoft YaHei", sans-serif;
```

### 行高

| 场景 | 行高 |
|------|------|
| 正文 | 1.6 |
| 表格行 | 约 44-46px (13px padding top/bottom) |

---

## 3. 色彩体系 (Colors)

### 中性色

| 变量名 | 值 | 用途 |
|--------|-----|------|
| --text-primary | #2b2b2b | 主要文字（标题、正文、表格内容、菜单选中） |
| --text-secondary | #6b6258 | 次要文字（描述、菜单项、按钮文字） |
| --text-muted | #b0a89a | 辅助文字（面包屑、统计卡标题、页脚） |
| --text-icon | #c4bdb0 | 图标色、Input placeholder |

### 品牌色

| 变量名 | 值 | 用途 |
|--------|-----|------|
| --brand-dark | #2b2b2b | 主按钮背景、激活页码、强调色 |
| --brand-dark-hover | #4d4d4d | 主按钮 hover (变浅) |

### 边框与分割线

| 变量名 | 值 | 用途 |
|--------|-----|------|
| --border-color | #e8e2d8 | 默认边框（输入框、按钮边框、分页项） |
| --border-hover | #d4cdc0 | 悬浮边框（按钮 hover、分页项 hover） |
| --border-light | #efeae2 | 浅边框（卡片、表格外框） |
| --divider | #f5f2ed | 分割线（表格行间、搜索栏底部分割） |

### 背景色

| 变量名 | 值 | 用途 |
|--------|-----|------|
| --bg-white | #ffffff | 纯白（表格列背景、输入框、激活页码文字） |
| --hover-bg | #f0ece6 | 悬浮背景（IconButton hover） |
| --hover-bg-light | #f5f2ed | 浅悬浮（次按钮 hover、分页项 hover） |
| --active-bg | #E4E0D8 | 选中/激活背景（菜单项选中、菜单 hover） |
| --bg-light | #faf8f5 | 备用浅背景（表头背景、行 hover） |

### 状态色

| 含义 | 主色 | 浅色背景 | 文字色 |
|------|------|----------|--------|
| 成功 / 通过 / 活跃 | #22c55e (--green) | #dcfce7 (--green-light) | #166534 (--green-text) |
| 失败 / 拒绝 / 异常 | #e74c3c (--red) | #fef2f2 (--red-light) | #e74c3c (--red-text) |
| 待处理 / 警告 | #f59e0b (--yellow) | #fef3c7 (--yellow-light) | #92400e (--yellow-text) |
| 信息 / 标识 | #3b6fdf (--blue) | #edf2ff (--blue-light) | #3b6fdf (--blue-text) |
| 默认 / 未分配 / 禁用 | #6b6258 (--gray) | #f5f2ed (--gray-light) | #b0a89a (--gray-text) |

---

## 4. 间距规范 (Spacing)

### 圆角

| 变量名 | 值 | 用途 |
|--------|-----|------|
| --radius-sm | 8px | 按钮、输入框、标签、分页项、菜单项、Select |
| --radius-md | 12px | 卡片、Modal、DataPanel 内容区、侧边栏菜单项 |
| --radius-lg | 16px | 预留，非标准大圆角场景 |

### 阴影

| 变量名 | 值 | 用途 |
|--------|-----|------|
| --shadow-sm | 0 1px 3px rgba(0,0,0,0.06) | 备用 |
| --shadow-md | 0 4px 12px rgba(0,0,0,0.08) | StatCard hover 抬起效果 |
| --shadow-lg | 0 8px 24px rgba(0,0,0,0.12) | 预留 |

### 页面级间距

| 位置 | 值 |
|------|-----|
| Page Header 上内边距 | 20px |
| Page Header 左右内边距 | 28px |
| Page Header ↔ Filter Bar | 16px |
| Filter Bar 左右内边距 | 28px |
| Filter Bar ↔ Content | 16px |
| Content 底部外间距 | 20px |
| Content 区容器边框 | 1px solid #efeae2 |
| Content 区圆角 | 12px (--radius-md) |
| 表格内单元格 padding | 12px 14px (表头) / 13px 14px (行) |
| 分页栏内边距 | 12px 18px |

---

## 5. 按钮体系 (Buttons)

### 尺寸规范

| 属性 | 值 |
|------|-----|
| 高度 | 34px |
| 圆角 | 8px (--radius-sm) |
| 字号 | 13px |
| 字重 | 500 |
| 过渡 | all 0.15s ease |

### 主按钮 (Primary)

| 状态 | 背景色 | 文字色 | 边框 |
|------|--------|--------|------|
| normal | #2b2b2b (--brand-dark) | #ffffff | #2b2b2b |
| hover | #4d4d4d (--brand-dark-hover) | #ffffff | #4d4d4d |
| disabled | antd 默认灰色 | — | — |

交互方向：hover 时背景 **变浅**。

### 次按钮 (Default)

| 状态 | 背景色 | 文字色 | 边框 |
|------|--------|--------|------|
| normal | transparent | #6b6258 (--text-secondary) | #e8e2d8 (--border-color) |
| hover | #f5f2ed (--hover-bg-light) | #6b6258 | #d4cdc0 (--border-hover) |
| disabled | antd 默认 | — | — |

交互方向：hover 时背景 **加深**（出现浅米色填充）。

### 文字按钮 (Text)

| 状态 | 背景色 | 文字色 |
|------|--------|--------|
| normal | transparent | #6b6258 |
| hover | #E4E0D8 (--active-bg) | #6b6258 |

### 表格操作按钮 (Table Action Buttons)

适用于表格「操作」列中的编辑、配置、删除等动作按钮。

**统一模式：**

```tsx
<Button type="link" size="small" icon={<EditOutlined />}>编辑</Button>
<Button type="link" size="small" danger icon={<DeleteOutlined />}>删除</Button>
```

**规范明细：**

| 属性 | 值 |
|------|-----|
| 按钮类型 | `link` |
| 尺寸 | `small` |
| 图标 | 通过 `icon` 属性传入，与文字间 gap 1px |
| 文字 | 中文动词（编辑、配置权限、查看、通过、退回等） |
| 非危险操作色 | #3d5a80（CSS 覆盖 `.ant-btn-link:not(.ant-btn-dangerous)`） |
| 危险操作颜色 | 保持 antd 默认红色（`danger` prop） |
| 间距 | gap 1px（CSS `gap` 属性，非 margin） |
| 按钮间分隔线 | `::before` 伪元素，1px 高，14px 高，#e8e2d8 色 |

**分隔线实现（CSS in global.css）：**

```css
.ant-table-tbody td .ant-space-item + .ant-space-item::before {
  content: '';
  width: 1px;
  height: 14px;
  background: var(--border-color);
  margin-right: 4px;
  flex-shrink: 0;
}
```

非表格场景使用 `<Space size={0} className="action-btn-group">` 包裹按钮组以应用相同分隔线。

---

## 6. 表格规范 (Table)

### 表头

| 属性 | 值 |
|------|-----|
| 字体大小 | 12px |
| 字重 | 600 |
| 颜色 | #8a8276 (--table-header-text) |
| 背景 | #faf8f5 (--table-header-bg) |
| 间距 | padding 12px 14px (首列 18px, 末列 18px) |
| 转换 | text-transform: uppercase |
| 底部边框 | 1px solid #efeae2 (--table-border) |

### 行

| 属性 | 值 |
|------|-----|
| 字体大小 | 13px |
| 颜色 | #2b2b2b (--text-primary) |
| 间距 | padding 13px 14px (首列 18px, 末列 18px) |
| 分割线 | border-bottom: 1px solid #f5f2ed (--table-row-divider) |
| hover 背景 | #faf8f5 (--table-hover-bg) |
| 末行 | 无底部边框 |

### 分页

| 属性 | 值 |
|------|-----|
| 内边距 | 12px 18px |
| 上边框 | 1px solid #efeae2 (--table-border) |
| 页码尺寸 | 30px x 30px |
| 页码圆角 | 8px (--radius-sm) |
| 页码边框 | #e8e2d8 (--border-color) |
| 页码文字 | 13px, #6b6258 |
| 当前页背景 | #2b2b2b (--brand-dark) |
| 当前页文字 | #ffffff |
| 分页项 hover | background #f5f2ed, border #d4cdc0 |
| 激活项 hover | 保持 #2b2b2b 背景，opacity 0.85 | 不变色 |

---

## 7. 标签体系 (Tag)

### 尺寸

| 属性 | 值 |
|------|-----|
| 圆角 | 6px |
| 字号 | 12px |
| 字重 | 500 |
| 内边距 | 2px 10px |
| 行高 | 20px |
| 无边框 | border: none |

### 色值映射

| 状态 | 背景色 | 文字色 | 对应 Tag 类型 |
|------|--------|--------|--------------|
| 成功 / 活跃 | #dcfce7 | #166534 | green |
| 失败 / 异常 | #fef2f2 | #e74c3c | red |
| 待处理 / 警告 | #fef3c7 | #92400e | yellow |
| 信息 / 角色标识 | #edf2ff | #3b6fdf | blue |
| 默认 / 未分配 | #f5f2ed | #b0a89a | gray |

---

## 8. 表单规范 (Form)

| 属性 | 值 |
|------|-----|
| 布局 | vertical (上下布局) |
| 标签字体 | 13px, 500 |
| 输入框高 | 34px |
| 输入框圆角 | 8px (--radius-sm) |
| 输入框边框 | 1px solid #e8e2d8 |
| 输入框 focus 边框 | #2b2b2b (--brand-dark) |
| 输入框 placeholder | #c4bdb0 (--text-icon) |
| 表单间距 (InputNumber addon) | addonAfter="单位" |
| Switch 状态文字 | checkedChildren / unCheckedChildren |
| 表单按钮间距 | gap 12px |

---

## 9. 筛选栏与工具栏规范 (Filters & Toolbar)

### 三个区域的职责分工

| 区域 | 位置 | 典型内容 | 对齐 |
|------|------|----------|------|
| filters | 标题下方筛选栏 | 搜索框(FilterSearch)、Select下拉、日期选择、查询/重置按钮 | 左对齐 |
| toolbarActions | 筛选栏右侧 | 导出、刷新等辅助操作 | 右对齐 |
| extra | 标题右侧 | 新建XX等主要操作 | 右对齐(标题行) |

### filters 区域规范

| 属性 | 值 |
|------|-----|
| 搜索框 | FilterSearch 组件，最大宽度 320px，左侧带搜索图标 |
| 查询按钮 | Default 次要按钮，icon 带 SearchOutlined |
| 重置按钮 | 白色背景次按钮 |
| 间隔 | 组件之间 gap 12px |

### toolbarActions 规范

| 属性 | 值 |
|------|-----|
| 按钮类型 | Default 次要按钮（非 Primary），34px 高 |
| 典型操作 | 导出(ExportOutlined)、刷新(ReloadOutlined) |
| 对齐 | 右对齐 |

### extra 规范

| 属性 | 值 |
|------|-----|
| 按钮类型 | Primary 主按钮 |
| 典型操作 | 新增XX |
| 位置 | 页面标题右侧 |

---

## 10. 交互规范 (Interaction)

### 悬浮 (hover) 规则

| 元素 | normal | hover | 方向 |
|------|--------|-------|------|
| 主按钮 (Primary) | #2b2b2b | #4d4d4d | **变浅** |
| 次按钮 (Default) | transparent | #f5f2ed | 显现填充 |
| 分页项 (非激活) | transparent | #f5f2ed / #d4cdc0 | 显现填充 |
| 表格行 | — | #faf8f5 | 显现浅背景 |
| 菜单项 | transparent | #E4E0D8 | 加深 |
| 菜单选中 | #E4E0D8 | #E4E0D8 | 不变(保持选中态) |
| 文字按钮 | transparent | #E4E0D8 | 显现填充 |
| Sidebar Toggle | transparent | #E4E0D8 | 显现填充 |

### 过渡时间

| 属性 | 值 |
|------|-----|
| 默认过渡 | all 0.15s ease |
| CSS 变量过渡 | 0.25s cubic-bezier(0.4, 0, 0.2, 1) |

---

## 11. CSS 变量参考

完整的 CSS 变量定义见 `admin/src/styles/theme.css`，全局覆盖见 `admin/src/styles/global.css`。

### 变量清单

| 分类 | 变量名 | 默认值 |
|------|--------|--------|
| **布局** | --sidebar-bg | #F5F3EF |
| | --sidebar-collapsed | 56px |
| | --sidebar-expanded | 224px |
| | --main-bg | #FFFFFF |
| | --content-bg | #FAFAF8 |
| **字体** | --font-family | -apple-system, "PingFang SC"... |
| | --text-primary | #2b2b2b |
| | --text-secondary | #6b6258 |
| | --text-muted | #b0a89a |
| | --text-icon | #c4bdb0 |
| **边框** | --border-color | #e8e2d8 |
| | --border-hover | #d4cdc0 |
| | --border-light | #efeae2 |
| | --divider | #f5f2ed |
| **背景** | --bg-white | #ffffff |
| | --bg-light | #faf8f5 |
| | --hover-bg | #f0ece6 |
| | --hover-bg-light | #f5f2ed |
| | --active-bg | #E4E0D8 |
| **品牌** | --brand-dark | #2b2b2b |
| | --brand-dark-hover | #4d4d4d |
| **状态** | --green (+light,+text) | #22c55e / #dcfce7 / #166534 |
| | --red (+light,+text) | #e74c3c / #fef2f2 / #e74c3c |
| | --yellow (+light,+text) | #f59e0b / #fef3c7 / #92400e |
| | --blue (+light,+text) | #3b6fdf / #edf2ff / #3b6fdf |
| | --gray (+light,+text) | #6b6258 / #f5f2ed / #b0a89a |
| **表格** | --table-header-bg | #faf8f5 |
| | --table-header-text | #8a8276 |
| | --table-border | #efeae2 |
| | --table-row-divider | #f5f2ed |
| | --table-hover-bg | #faf8f5 |
| **圆角** | --radius-sm/-md/-lg | 8px / 12px / 16px |
| **阴影** | --shadow-sm/-md/-lg | 三级阴影 |
| **过渡** | --transition | 0.25s cubic-bezier(0.4,0,0.2,1) |
